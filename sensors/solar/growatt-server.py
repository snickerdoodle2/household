import growattServer
from jproperties import Properties

configs = Properties()
with open('credentials.properties', 'rb') as config_file:
    configs.load(config_file)


api = growattServer.GrowattApi()
login_response = api.login(configs.get("USERNAME").data, configs.get("PASSWORD").data)
# Get a list of growatt plants.
print(api.plant_list(login_response['user']['id']))
