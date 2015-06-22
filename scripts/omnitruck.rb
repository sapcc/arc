#!/usr/bin/env ruby
require 'json'
require 'webrick'
server = WEBrick::HTTPServer.new(:Port => ENV.fetch('PORT', 4567), :DocumentRoot => File.expand_path('../../tmp', __FILE__) )

msi = {
  relpath: '/chef-client-12.3.0-1.msi',
  md5: '842a3c49b52bdd9f838be0c6d63fd68c',
  sha256: '13fd1a79a91be692a63fba0c77e7bebe9a7173047a4dd22e7b5eae74393804d4'
}

rpm = {
  relpath: '/chef-12.3.0-1.el6.x86_64.rpm',
  md5: 'c19fefcb3d033107e9fbdb3839312584',
  sha256: '4b7c846a9ad93564cc203a5ac99890431f7d6ad159c424aa89827fd772c9881d'
}

deb = {
  relpath: '/chef_12.3.0-1_amd64.deb',
  md5: 'd8421c9b3010deb03e713ada00387e8a',
  sha256: 'e06eb748e44d0a323f4334aececdf3c2c74d2f97323678ad3a43c33ac32b4f81'
}

server.mount_proc("/metadata") do |req,res|

  result = case req.query["p"]
           when 'windows' then msi.clone
           when 'redhat','suse' then rpm.clone
           when 'ubuntu' then deb.clone
  end
  if result 
    result[:url]='http://'+req.header['host'].first + result[:relpath]
    res.body << JSON.pretty_generate(result)
  else
    res.status=404
    res.body<< "404 not found\n"
  end
end
server.start
