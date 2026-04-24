require "net/http"
require "json"
require "uri"
require "yaml"

class DiscogsService
  BASE_URL = "https://api.discogs.com".freeze
  USER_AGENT = "SerendipityMusicPlayer/1.0".freeze
  MAX_RETRIES = 5
  GENRES_PATH = Rails.root.join("config/discogs_genres.yml").freeze

  class ApiError < StandardError; end

  def self.genres
    @genres ||= YAML.load_file(GENRES_PATH)
  end

  def self.random_playable_track(genres: nil)
    new.random_playable_track(genres: genres)
  end

  def random_playable_track(genres: nil)
    MAX_RETRIES.times do
      begin
        release = fetch_random_release(genre: genres&.sample)
        next unless release

        track = extract_track_from_release(release)
        return track if track
      rescue ApiError => e
        Rails.logger.warn("Discogs API error: #{e.message}")
        next
      end
    end

    nil
  end

  private

  def fetch_random_release(genre: nil)
    random_page = rand(1..10_000)
    query = "type=release&per_page=1&page=#{random_page}"
    query += "&genre=#{URI.encode_www_form_component(genre)}" if genre.present?
    search_url = URI("#{BASE_URL}/database/search?#{query}")

    search_response = make_request(search_url)
    return nil unless search_response["results"]&.any?

    release_id = search_response["results"].first["id"]
    release_url = URI("#{BASE_URL}/releases/#{release_id}")

    make_request(release_url)
  end

  def extract_track_from_release(release)
    videos = release["videos"]
    return nil unless videos&.any?

    video = videos.find { |v| extract_youtube_id(v["uri"]).present? }
    return nil unless video

    video_id = extract_youtube_id(video["uri"])
    return nil unless video_id

    {
      title: video["title"] || "Unknown Track",
      artist: release.dig("artists", 0, "name") || "Unknown Artist",
      video_id: video_id
    }
  end

  def extract_youtube_id(uri)
    return nil if uri.blank?

    uri = uri.to_s

    if uri.include?("youtube.com/watch")
      query = URI.parse(uri).query
      return nil unless query

      params = URI.decode_www_form(query).to_h
      params["v"]
    elsif uri.include?("youtu.be/")
      path = URI.parse(uri).path
      path&.sub("/", "")
    else
      nil
    end
  rescue URI::InvalidURIError
    nil
  end

  def make_request(uri)
    http = Net::HTTP.new(uri.host, uri.port)
    http.use_ssl = true
    http.open_timeout = 10
    http.read_timeout = 10

    request = Net::HTTP::Get.new(uri)
    request["User-Agent"] = USER_AGENT

    response = http.request(request)

    case response.code.to_i
    when 200
      JSON.parse(response.body)
    when 429
      raise ApiError, "Rate limited by Discogs API"
    else
      raise ApiError, "Discogs API returned #{response.code}"
    end
  rescue JSON::ParserError => e
    raise ApiError, "Failed to parse Discogs response: #{e.message}"
  rescue SocketError, Net::OpenTimeout, Net::ReadTimeout => e
    raise ApiError, "Network error: #{e.message}"
  end
end
