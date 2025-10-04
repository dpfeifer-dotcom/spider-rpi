from evdev import InputDevice, categorize, ecodes, list_devices
import controllers

# Nyisd meg a 0. eszközt
def find_gamepad():
    devices = [InputDevice(path) for path in list_devices()]
    for dev in devices:
        capabilities = dev.capabilities()
        print(capabilities.get(ecodes.EV_KEY, []))
        if ecodes.EV_ABS in capabilities or ecodes.EV_KEY in capabilities:
            if any(code in capabilities.get(ecodes.EV_KEY, []) for code in [
                ecodes.BTN_GAMEPAD,
                ecodes.BTN_SOUTH, ecodes.BTN_EAST, ecodes.BTN_NORTH, ecodes.BTN_WEST
            ]):
                return dev
    return None

joystick = find_gamepad()
print(f"Using device: {joystick.name} ({joystick.path})")

# Xbox gombok kódjai az evdev-ben
BUTTON_MAP = {
    ecodes.BTN_SOUTH: 'A',
    ecodes.BTN_EAST: 'B',
    ecodes.BTN_NORTH: 'X',
    ecodes.BTN_WEST: 'Y',
    ecodes.BTN_TR: 'BTN-TR',
    ecodes.BTN_TL: 'BTN-TL',
    ecodes.ABS_HAT0X: 'BTN_X',
    ecodes.ABS_HAT0Y: 'BTN_Y'
    
}

for event in joystick.read_loop():
    if event.type == ecodes.EV_KEY:
        button_name = BUTTON_MAP.get(event.code)
        if button_name == 'Y' and event.value:
            controllers.switch_light("police_warn")
        if button_name == 'B' and event.value:
            controllers.switch_light("police_calm")
        if button_name == 'X' and event.value:
            controllers.switch_light("warning")
        if button_name == 'A' and event.value:
            controllers.switch_light("none")
        if button_name == 'BTN_TR' and event.value:
            controllers.switch_bulb()
        else:
            btnevent = categorize(event)
            print(f"Button {absevent.event.code} value {absevent.event.value}")

    elif event.type == ecodes.EV_ABS:
        absevent = categorize(event)
        print(f"Axis {absevent.event.code} value {absevent.event.value}")