installer = ::File.join Chef::Config[:file_cache_path], node['golang']['installer_filename']
remote_file installer do
  source node['golang']['installer_url']
end

windows_package "go" do
  source installer
end

path = File.dirname node['golang']['dev_pkg']
pkg  = File.basename node['golang']['dev_pkg']

directory "C:/gocode/src/#{path}" do
  recursive true
end
execute "mklink /J C:\\gocode\\src\\#{node['golang']['dev_pkg'].gsub('/',"\\")} C:\\vagrant" do
  not_if {File.exists? "C:/gocode/src/#{node['golang']['dev_pkg']}"}
end
env "GOPATH" do
  value "C:\\vagrant\\Godeps\\_workspace;C:\\gocode"
end

