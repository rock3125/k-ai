

#
# process the semantics and set parentage based on types
#

person = ["person", "mother", "male_parent", "pal", "comrade", "faculty_member", "academic", "officer", "lawman", "man", "woman", "old_person"]
location = ["location", "workplace", "room", "area", "dwell", "business", "home", "domicile"]
vehicle = ["vehicle"]
vehicle_negative = ["craft", "booster"]
aircraft = ["aircraft"]
aircraft_negative = ["plane"]
container = ["container"]
plant = ["vascular_plant", "fungus", "crop", "pot_plant", "perennial", "hygrophyte", "monocarpic_plant",
         "monocarpous_plant", "sporophyte", "shrub"]
animal = ["canine", "feline", "cattle", "eutherian", "eutherian_mammal", "saltwater_fish", "freshwater_fish",
          "elasmobranch", "food_fish", "salmonid", "big_cat", "spider", "ant", "insect", "dipteran", "dipteron",
          "insect", "arthropod", "worm", "hoofed_mammal","even-toed_ungulate", "artiodactyl", "artiodactyl_mammal",
          "dinocerate", "odd-toed_ungulate", "perissodactyl", "perissodactyl_mammal", "dinocerate", "even-toed_ungulate",
          "artiodactyl", "artiodactyl_mammal", "ungulate", "odd-toed_ungulate", "perissodactyl", "perissodactyl_mammal", "ungulate",
          "thrips", "louse"]

# load existing semantics
def get_semantics(data_list, negative, num_levels):
    set1 = {}
    exists = {}
    forbidden = {}
    list = []
    for item in data_list:
        set1[item] = 1
        exists[item] = 1
        list.append(item)
    if negative is not None:
        for item in negative:
            forbidden[item] = 1

    for i in range(num_levels):
        with open('/opt/kai/data/wordnet/wordnet-3.1-relationship-graph.txt') as reader:
            for line in reader:
                line = line.strip()
                if len(line) > 0 and not line.startswith("#"):
                    parts = line.split("|")
                    if len(parts) == 3:
                        if ":n" in parts[0] and ":n" in parts[2] and ("2" in parts[1]):
                            word1 = parts[0].split(":")[0]
                            word2 = parts[2].split(":")[0]
                            if word2 in set1 and word2 not in forbidden and not word1 in exists:
                                list.append(word1)
                                exists[word1] = 1

        # repeat?
        for item in list:
            set1[item] = 1

    return list

def display(set, semantic):
    with open(semantic + ".txt", 'w') as writer:
        set.sort()
        for item in set:
            writer.write(item + ":" + semantic + "\n")

# people_set = get_semantics(person, None, 1)
# display(people_set, "person")
#
# location_set = get_semantics(location, None, 1)
# display(location_set, "location")
#
# vehicle_set = get_semantics(vehicle, vehicle_negative, 2)
# display(vehicle_set, "vehicle")
#
# aircraft_set = get_semantics(aircraft, aircraft_negative, 3)
# display(aircraft_set, "aircraft")
#
# container_set = get_semantics(container, None, 1)
# display(container_set, "container")
#
# plant_set = get_semantics(plant, None, 1)
# display(plant_set, "plant")

animal_set = get_semantics(animal, None, 1)
display(animal_set, "animal")
