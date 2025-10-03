#!/usr/bin/env python3
import pygame
import sys
import time

def init_joysticks():
    """Inicializálja a joystick eszközöket és visszaadja a listát."""
    print("START")
    pygame.joystick.init()  # csak joystick modult inicializáljuk
    joysticks = []
    for i in range(pygame.joystick.get_count()):
        js = pygame.joystick.Joystick(i)
        js.init()
        joysticks.append(js)
        print(f"Észlelt eszköz: {js.get_name()} (axes: {js.get_numaxes()}, buttons: {js.get_numbuttons()}, hats: {js.get_numhats()})")
    return joysticks

def main():
    joysticks = init_joysticks()

    if not joysticks:
        print("Nincs joystick/gamepad csatlakoztatva.")
        # várjunk egy kicsit, hogy a konténer ne zárjon azonnal
        time.sleep(5)
        sys.exit(0)

    running = True
    try:
        while running:
            # Csak a joystick eseményeket olvassuk, ablak nélkül
            for js in joysticks:
                for i in range(js.get_numbuttons()):
                    if js.get_button(i):
                        print(f"Gomb lenyomva: joy={js.get_id()}, button={i}")

                for a in range(js.get_numaxes()):
                    val = js.get_axis(a)
                    if abs(val) > 0.01:  # kis zajt figyelmen kívül hagyjuk
                        print(f"Tengely mozgás: joy={js.get_id()}, axis={a}, value={val:.3f}")

                for h in range(js.get_numhats()):
                    val = js.get_hat(h)
                    if val != (0, 0):
                        print(f"Hat switch: joy={js.get_id()}, hat={h}, value={val}")

            time.sleep(0.01)  # kis delay, ne pörögjön a CPU
    except KeyboardInterrupt:
        print("Kilépés...")

if __name__ == "__main__":
    main()
