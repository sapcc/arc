require 'net/http'
require 'json'
uri = URI(ARGV.first)

http = Net::HTTP.new(uri.host, uri.port)

request = Net::HTTP::Post.new(uri.path, {'Content-Type' =>'application/json'})
request.body={app_id:"arc", app_version:"0.1.0", tags:{arch:"amd64", channel:"stable", os:"windows"}}.to_json
response = http.request(request)
puts response.code + " " + response.msg
puts JSON.pretty_generate(JSON.parse(response.body)) if response.body
