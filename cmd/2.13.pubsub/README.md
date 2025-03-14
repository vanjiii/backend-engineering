podman run --rm --hostname rabbit -p 5672:5672 -p 15671:15672 rabbitmq:4.1-rc-management
# 5672 - access protocol
# 15672 - dashboard port
    note: the image needs to have `*-management` in the name to have this plugin installed
    guest:guest
