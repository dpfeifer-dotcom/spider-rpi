from evdev import InputDevice, categorize, ecodes

# Nyisd meg a 0. eszközt
joystick = InputDevice('/dev/input/event0')

print(f"Using device: {joystick.name} ({joystick.path})")

# Végtelen olvasó loop
for event in joystick.read_loop():
    if event.type == ecodes.BTN_SOUTH:
        # Gomb esemény
        state = 'pressed' if event.value else 'released'
        print(f"Button {event.code} {state}")
    if event.type == ecodes.BTN_NORTH:
        # Gomb esemény
        state = 'pressed' if event.value else 'released'
        print(f"Button {event.code} {state}")
    if event.type == ecodes.BTN_WEST:
        # Gomb esemény
        state = 'pressed' if event.value else 'released'
        print(f"Button {event.code} {state}")
    if event.type == ecodes.BTN_EAST:
        # Gomb esemény
        state = 'pressed' if event.value else 'released'
        print(f"Button {event.code} {state}")
    elif event.type == ecodes.EV_ABS:
        # Tengely esemény
        absevent = categorize(event)
        print(f"Axis {absevent.event.code} value {absevent.event.value}")




