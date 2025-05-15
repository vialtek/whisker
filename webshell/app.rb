require 'rubygems'
require 'sinatra'
require 'sinatra/json'

# Static Jobs as examples
available_jobs = [
    {
      guid: SecureRandom.uuid.gsub('-','').downcase,
      workflow: 'hello',
      dataset: 'beer',
      status: :waiting
    },
    {
      guid: SecureRandom.uuid.gsub('-','').downcase,
      workflow: 'hello',
      dataset: 'beer',
      status: :waiting
    }
  ]

get '/' do
  erb :index
end

post '/heartbeat' do
  params = JSON.parse(request.body.read, symbolize_names: true)
  puts params

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

patch '/jobs/:guid/accept' do
  job = available_jobs.detect { it[:guid] == params['guid']}
  return json status: 'not_found' if job == nil

  job[:status] = :in_progress
  return json status: 'ok'
end

patch '/jobs/:guid/finished' do
  job = available_jobs.detect { it[:guid] == params['guid']}
  return json status: 'not_found' if job == nil

  job[:status] = :finished
  return json status: 'ok'
end

patch '/jobs/:guid/failed' do
  job = available_jobs.detect { it[:guid] == params['guid']}
  return json status: 'not_found' if job == nil

  job[:status] = :failed
  return json status: 'ok'
end

post '/jobs/:guid/output_log' do
  job = available_jobs.detect { it[:guid] == params['guid']}
  return json status: 'not_found' if job == nil

  params = JSON.parse(request.body.read, symbolize_names: true)
  job[:output_log] = params[:output_log]

  return json status: 'ok'
end
