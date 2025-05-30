#!/usr/bin/env python3
"""
CivIdleCli - A Command Line Civilization Builder Game
"""
import sys
import os
from rich.console import Console

# Add the project directory to the path so we can import our modules
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from game.engine import GameEngine
from ui.display import Display

def main():
    console = Console()
    console.print("[bold green]CivIdleCli - A Command Line Civilization Builder[/bold green]", justify="center")
    
    # Initialize the game
    display = Display(console)
    game = GameEngine(display)
    
    # Show intro and start the game
    display.show_intro()
    game.start()

if __name__ == "__main__":
    main()
