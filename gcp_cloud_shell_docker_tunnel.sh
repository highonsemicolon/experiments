#!/bin/bash

# Execute remote commands
gcloud cloud-shell ssh --authorize-session --command='bash -c "
    echo Hello from Cloud Shell!

    # Enable Docker API on localhost:2375 if not enabled
    if ! sudo grep -q \"tcp://127.0.0.1:2375\" /etc/docker/daemon.json; then
        echo Configuring Docker API to listen on localhost:2375...

        echo Running additional commands...

        cat <<EOF | sudo tee /etc/docker/daemon.json > /dev/null 
{
   \"hosts\": [\"tcp://127.0.0.1:2375\", \"unix:///var/run/docker.sock\"]
}
EOF

        sudo service docker restart
    fi
"'

gcloud cloud-shell ssh --authorize-session --ssh-flag="-L 2375:localhost:2375"

