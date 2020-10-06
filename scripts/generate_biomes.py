from ctypes import *
from datetime import datetime
import os.path as path
import urllib.request
import json

# TODO: switch to using go:generate

# This has all of the data we need
import stringcase as stringcase

with urllib.request.urlopen("https://raw.githubusercontent.com/PrismarineJS/minecraft-data/master/data/pc/1.15.2/biomes.json") as f:
    biomes = json.load(f)


def figure_biome_type(name):
    # Patterns
    if 'ocean' in name:
        return 'TypeOcean'
    elif 'plains' in name:
        return 'TypePlains'
    elif 'hills' in name:
        return 'TypeHills'
    elif 'desert' in name:
        return 'TypeDesert'
    elif 'forest' in name:
        return 'TypeForest'
    elif 'taiga' in name:
        return 'TypeTaiga'
    elif 'swamp' in name:
        return 'TypeSwamp'
    elif 'river' in name:
        return 'TypeRiver'
    elif 'nether' in name:
        return 'TypeNether'
    elif 'end' in name:
        return 'TypeSky'
    elif 'snow' in name:
        return 'TypeSnow'
    elif 'beach' in name:
        return 'TypeBeach'
    elif 'mushroom' in name:
        return 'TypeMushroomIsland'
    elif 'jungle' in name:
        return 'TypeJungle'
    elif 'savanna' in name:
        return 'TypeSavanna'
    elif 'badlands' in name:
        return 'TypeMesa'

    elif 'stone_shore' == name:
        return 'TypeStoneBeach'

    elif 'the_void' == name:
        return 'TypeVoid'

    elif 'ice_spikes' == name:
        return 'TypeSnow'

    # must be last so it will get the most correct of all
    elif 'mountain' in name:
        return 'TypeHills'

    else:
        assert False, f"Could not figure out type for {name}"


print("// Code generated by scripts/generate_biomes.go; DO NOT EDIT.")
print("// This file was generated by robots at")
print("// " + str(datetime.now()))
print()

print("package biome")
print()

# Because the world gen uses biomes in a very special
# way we are going to create variable for each biome
# type for ease of use

for biome in biomes:
    name = stringcase.pascalcase(biome['name'])
    print(f"var {name} = &Biome{{")
    print(f"\tId: {biome['id']},")
    print(f"\tName: \"{biome['name']}\",")
    print(f"\tType: {figure_biome_type(biome['name'])},")
    print(f"\tRainfall: {biome['rainfall']},")
    print(f"\tTemperature: {biome['temperature']},")
    print("}")

print()

# Biome list by id
print("var biomes = [...]*Biome{")
for biome in biomes:
    name = stringcase.pascalcase(biome['name'])
    print(f"\t{biome['id']}: {name},")
print("}")
print("")
