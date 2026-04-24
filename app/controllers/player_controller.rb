class PlayerController < ApplicationController
  def index
  end

  def random_track
    genres = params[:genres].presence
    genres = Array(genres).reject(&:blank?) if genres

    track = DiscogsService.random_playable_track(genres: genres)

    if track
      render json: track
    else
      render json: { error: "No playable track found. Please try again." }, status: :not_found
    end
  end
end
