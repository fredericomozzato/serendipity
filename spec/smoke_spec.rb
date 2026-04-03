# frozen_string_literal: true

require 'rails_helper'

RSpec.describe "RSpec wiring" do
  it "runs successfully" do
    expect(true).to be true
  end

  it "has Rails loaded" do
    expect(defined?(Rails)).to eq("constant")
  end
end
