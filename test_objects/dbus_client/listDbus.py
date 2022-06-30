import dbus
for service in dbus.SystemBus().list_names():
        print(service)

for service in dbus.SessionBus().list_names():
        print(service)
