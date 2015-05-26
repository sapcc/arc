default['golang']['version']= '1.4.2'
default['golang']['installer_filename'] = "go#{node['golang']['version']}.windows-amd64.msi"
default['golang']['installer_url']= "https://storage.googleapis.com/golang/#{node['golang']['installer_filename']}"
default['golang']['dev_pkg']='gitHub.***REMOVED***/monsoon/arc'


