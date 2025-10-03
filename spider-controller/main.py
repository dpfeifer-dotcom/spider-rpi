#!/usr/bin/env python3
"""
gamepad_detector_pygame.py

Egyszerű Python alkalmazás, ami pygame segítségével érzékeli a csatlakoztatott USB gamepadokat, és kiírja a gomb- és tengelyeseményeket.

Követelmények:
  - python3
  - pip install pygame

Futtatás: python3 gamepad_detector_pygame.py

Megjegyzés: a pygame cross-platform, így Windows/Mac/Linux alatt is működik.
"""

import pygame
import sys


def init_joysticks():
    """Inicializálja a joystick eszközöket és visszaadja a listát."""
    pygame.joystick.init()
    joysticks = []
    for i in range(pygame.joystick.get_count()):
        js = pygame.joystick.Joystick(i)
        js.init()
        joysticks.append(js)
        print(f"Észlelt eszköz: {js.get_name()} (axes: {js.get_numaxes()}, buttons: {js.get_numbuttons()}, hats: {js.get_numhats()})")
    return joysticks


def main():
    clock = pygame.time.Clock()
    joysticks = init_joysticks()

    if not joysticks:
        print("Nincs joystick/gamepad csatlakoztatva.")

    running = True
    while running:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                running = False

            # Joystick események
            elif event.type == pygame.JOYBUTTONDOWN:
                print(f"Gomb lenyomva: joy={event.joy}, button={event.button}")

            elif event.type == pygame.JOYBUTTONUP:
                print(f"Gomb felengedve: joy={event.joy}, button={event.button}")

            elif event.type == pygame.JOYAXISMOTION:
                print(f"Tengely mozgás: joy={event.joy}, axis={event.axis}, value={event.value:.3f}")

            elif event.type == pygame.JOYHATMOTION:
                print(f"Hat switch: joy={event.joy}, value={event.value}")

        # FPS limit, hogy ne pörögjön feleslegesen
        clock.tick(60)

    pygame.quit()
    sys.exit()


if __name__ == "__main__":
    main()
