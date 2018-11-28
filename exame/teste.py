import random
import time
import numpy as np
from copy import deepcopy


#mapa_original = np.loadtxt("map.txt", dtype='str', delimiter=',')
#mapa = np.loadtxt("map.txt", dtype='str', delimiter=',')
mapaChar = np.loadtxt("map1.txt", comments="#", delimiter=",", unpack=False)
print(mapaChar)