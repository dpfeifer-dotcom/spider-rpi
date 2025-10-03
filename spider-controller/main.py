from evdev import InputDevice, categorize, ecodes

# Nyisd meg a 0. eszközt
joystick = InputDevice('/dev/input/event0')

print(f"Using device: {joystick.name} ({joystick.path})")

# Xbox gombok kódjai az evdev-ben
BUTTON_MAP = {
    ecodes.BTN_SOUTH: 'A',
    ecodes.BTN_EAST: 'B',
    ecodes.BTN_NORTH: 'X',
    ecodes.BTN_WEST: 'Y',
}

# Végtelen olvasó loop
for event in joystick.read_loop():
    if event.type == ecodes.EV_KEY:
        button_name = BUTTON_MAP.get(event.code)
        if button_name:
            state = 'pressed' if event.value else 'released'
            print(f"Button {button_name} {state}")
    elif event.type == ecodes.EV_ABS:
        absevent = categorize(event)
        print(f"Axis {absevent.event.code} value {absevent.event.value}")