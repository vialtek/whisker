require 'rubygems'
require 'sinatra'
require 'sinatra/json'

get '/' do
  erb :index
end

post '/heartbeat' do
  json status: 'ok'
end

get '/jobs' do
  available_jobs = [
    {
      guid: '1',
      workflow: 'hello',
      dataset: ''
    }
  ]

  json available_jobs
end
