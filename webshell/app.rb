require 'rubygems'
require 'sinatra'
require 'sinatra/json'

seen_nodes = []
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
      workflow: 'call-recipe',
      dataset: 'beer',
      status: :waiting,
      params: {
        targetVariable: 'Extract Concentration',
        learningRate: 0.001

      }
    },
  ]

get '/' do
  erb :index
end

get '/nodes' do
  json seen_nodes
end

post '/heartbeat' do
  params = JSON.parse(request.body.read, symbolize_names: true)
  puts params

  node = seen_nodes.detect { it[:node_name] == params[:node_name] }
  if node.nil?
    node = {node_name: params[:node_name]}
    seen_nodes.push node
  end

  node[:last_seen] = DateTime.now
  node[:busy] = params[:busy] == 'true'
  node[:datasets] = params[:datasets]
  node[:workflows] = params[:workflows]
  node[:jobs_processed] = params[:jobs_processed]
  node[:uptime] = params[:uptime]

  json status: 'ok'
end

get '/jobs' do
  json available_jobs.select { it[:status] == :waiting }
end

get '/jobs/all' do
  json available_jobs
end

get '/jobs/refresh' do
  available_jobs.each { it[:status] = :waiting }
  json available_jobs
end

get '/jobs/:guid' do
  job = available_jobs.detect { it[:guid] == params['guid']}
  return json status: 'not_found' if job.nil?

  json job
end

post '/jobs/:guid/output_log' do
  job = available_jobs.detect { it[:guid] == params['guid']}
  return json status: 'not_found' if job.nil?

  params = JSON.parse(request.body.read, symbolize_names: true)
  job[:output_log] = params[:output_log]

  return json status: 'ok'
end

AVAILABLE_JOB_STATES = [:accept, :finished, :failed]
post '/jobs/:guid/:state' do
  state = params[:state].to_sym
  return json status: 'state_incorect' unless AVAILABLE_JOB_STATES.include? state

  job = available_jobs.detect { it[:guid] == params['guid']}
  return json status: 'not_found' if job.nil?

  state = :in_progress if state == :accept
  job[:status] = state

  return json status: 'ok'
end
