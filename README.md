# Photobooth

## Install
Flash an SD card with Raspbian lite and boot it up with an ethernet connection.

Do the first time setup:
- Login with the `pi / raspberry` account
- `sudo raspi-config`
    - Set your locale / timezone / keyboard layout accordingly
    - Interface options: enable SSH
    - Change the raspberry account password
    - Reboot 

Clone the repository:
```
$ git clone https://github.com/oxodao/photobooth-v3.git
$ cd photobooth-v3/ansible
```

Fill the inventory correctly:
```
$ nvim inventories/hosts
photobooth ansible_user=raspberry ansible_password=[[ YOUR RPI PASSWORD ]] ansible_host=[[ YOUR RPI ADDRESS ]] ansible_port=22
```

Change the password for the admin app and optionally the other config in `ansible/inventories/group_vars/all.yml`

Once done you can process with the ansible script:
```sh
$ ansible-galaxy install -r requirements.yaml
$ ansible-playbook -i inventories/hosts setup.yaml
```

Reboot your Pi and you should be good to go!