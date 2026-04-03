require 'spec_helper'
ENV['RAILS_ENV'] ||= 'test'
require_relative '../config/environment'
require 'rspec/rails'
require 'vcr'
require 'webmock/rspec'

VCR.configure do |c|
  c.cassette_library_dir = 'spec/cassettes'
  c.hook_into :webmock
  c.ignore_localhost = true
  c.configure_rspec_metadata!
  c.default_cassette_options = { record: :new_episodes }
end

WebMock.disable_net_connect!(allow_localhost: true)

# V1: No database — prevent ActiveRecord::TestFixtures from attempting PostgreSQL
# connections. MinitestLifecycleAdapter's around hook calls before_setup/after_teardown
# on each example instance; these overrides short-circuit both.
module NoDbFixtureSetup
  def before_setup
    # V1: No-op — skip fixture setup, no DB available
  end

  def after_teardown
    # V1: No-op — skip fixture teardown, @fixture_connection_pools is nil without DB
  end
end

RSpec.configure do |config|
  config.use_active_record = false
  config.use_transactional_fixtures = false
  config.infer_spec_type_from_file_location!
  config.filter_rails_from_backtrace!

  # Include after rspec-rails includes FixtureSupport so this override wins
  config.include NoDbFixtureSetup
end
