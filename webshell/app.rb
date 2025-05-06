require 'rubygems'
require 'sinatra'
require 'sinatra/json'

# Static Jobs as examples
available_jobs = [
    {
      guid: SecureRandom.uuid.gsub('-','').downcase,
      workflow: 'hello',
      dataset: '',
      status: :waiting
    },
    {
      guid: SecureRandom.uuid.gsub('-','').downcase,
      workflow: 'hello',
      dataset: '',
      status: :waiting
    }
  ]

get '/' do
  erb :index
end

post '/heartbeat' do
  json status: 'ok'
end

get '/jobs' do
  json available_jobs.select { it[:status] == :waiting }
end

get '/jobs/all' do
  json available_jobs
end

get '/jobs/:guid' do
  job = available_jobs.detect { it[:guid] == params['guid']}
  return json status: 'not_found' if job == nil

  json job
end

get '/jobs/:guid/accept' do
  job = available_jobs.detect { it[:guid] == params['guid']}
  return json status: 'not_found' if job == nil

  job[:status] = :in_progress
  return json status: 'ok'
end
