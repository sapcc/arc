ENV['VAGRANT_SERVER_URL'] = "http://image-store.***REMOVED***:8089"

Vagrant.configure(VAGRANTFILE_API_VERSION = "2") do |config|

  config.vm.define :win do |config|
    config.vm.box = 'monsoon/win2008r2'
    config.vm.provider :vmware_fusion do |v, override|
      v.gui = true
    end
    config.vm.provision :chef_solo do |chef|
      chef.log_level = :info
      chef.cookbooks_path = %w{chef/cookbooks chef/site-cookbooks}
      chef.add_recipe "sap-proxy::default"
      chef.add_recipe "sap-ssl::default"
      chef.add_recipe "golang::default"
    end

    sub_packages = Dir.glob("./**/*.go").reject {|p|p =~ /Godep/}.map {|p| File.dirname(p) }.uniq 
    config.vm.provision :shell, inline: "cd C:\\vagrant; go test -v #{sub_packages.join(" ")}"
  end

end
