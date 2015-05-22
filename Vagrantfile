ENV['VAGRANT_SERVER_URL'] = "http://image-store.***REMOVED***:8089"

Vagrant.configure(VAGRANTFILE_API_VERSION = "2") do |config|

  config.vm.define :win do |config|
    config.vm.box = 'monsoon/win2008r2'
    config.vm.provider :vmware_fusion do |v, override|
      v.gui = true
    end
    config.vm.provision :chef_solo do |chef|
      chef.log_level = :debug
      chef.cookbooks_path = %w{chef/cookbooks}
      chef.add_recipe "sap-proxy::default"
      chef.add_recipe "dotfiles-fabus"
    end
  end

end
