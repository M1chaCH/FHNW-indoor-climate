import proto
import util

config_dictionary = {}

def set(name, default_value, type_id): 
    global config_dictionary
    config_dictionary[name] = {
        "name": name,
        "value": default_value,
        "type": type_id,
    }

def get(name):
    global config_dictionary
    return config_dictionary[name].get("value")

def get_or_set(name, default_value, type_id):
    global config_dictionary
    entry = config_dictionary.get(name)
    if entry == None:
        config_dictionary[name] = {
            "name": name,
            "value": default_value,
            "type": type_id,
        }
        return default_value
    return entry.get("value")


def get_proto_options():
    global config_dictionary
    options = []

    for config in config_dictionary.values():
        type = config.get("type")
        if type == proto.CONFIG_OPTION_TYPE_STRING:
            options.append(create_string_option(config.get("name"), config.get("value")))
        elif type == proto.CONFIG_OPTION_TYPE_INT32:
            options.append(create_int_option(config.get("name"), config.get("value")))
        elif type == proto.CONFIG_OPTION_TYPE_DOUBLE:
            options.append(create_double_option(config.get("name"), config.get("value")))
        elif type == proto.CONFIG_OPTION_TYPE_BOOL:
            options.append(create_flag_option(config.get("name"), config.get("value")))
        else:
            print("Skipping unknown option: ", config)

    return options

def create_string_option(name, value):
    return proto.ConfigOption(config_name=name, config_type=proto.CONFIG_OPTION_TYPE_STRING, string_value=value)

def create_int_option(name, value):
    return proto.ConfigOption(config_name=name, config_type=proto.CONFIG_OPTION_TYPE_INT32, int32_value=value)

def create_double_option(name, value):
    return proto.ConfigOption(config_name=name, config_type=proto.CONFIG_OPTION_TYPE_DOUBLE, double_value=value)

def create_flag_option(name, value):
    return proto.ConfigOption(config_name=name, config_type=proto.CONFIG_OPTION_TYPE_BOOL, flag_value=value)


def set_from_proto(remote_config):
    print("setting config from remote data", remote_config, remote_config.device_id, util.get_device_id())
    if util.get_device_id() != remote_config.device_id:
        print("got remote config without matching device id, ignoring")
        return False

    for option in remote_config.options:
        if option.config_type == proto.CONFIG_OPTION_TYPE_STRING or option.config_type == None: # None -> 0 ...
            set(option.config_name, option.string_value, proto.CONFIG_OPTION_TYPE_STRING)
        elif option.config_type == proto.CONFIG_OPTION_TYPE_DOUBLE:
            set(option.config_name, option.double_value, option.config_type)
        elif option.config_type == proto.CONFIG_OPTION_TYPE_BOOL:
            set(option.config_name, option.flag_value, option.config_type)
        elif option.config_type == proto.CONFIG_OPTION_TYPE_INT32:
            set(option.config_name, option.int32_value, option.config_type)
        else:
            print("got invalid option, skipping", option)
    
    return True
