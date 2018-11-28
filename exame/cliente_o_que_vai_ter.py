import random
import time
import numpy as np
from copy import deepcopy
from os import system

# Le arquivo txt e salva em uma matriz de char
n_linhas = 40
n_colunas = 120
dim = (n_linhas, n_colunas)
mapa_original = np.chararray(dim)
with open("map.txt") as f:
    i = 0
    j = 0
    while True:
        c = f.read(1)
        if not c:
            break
        if c != '\n':
            mapa_original[i][j] = c
            j = j + 1
        else:
            i = i + 1
            j = 0
mapa = deepcopy(mapa_original)

# Movimento aleatório do usuário
def movimenta(p):
    x = p[0]
    y = p[1]
    x_antes = x
    y_antes = y

    falhou = False
    while (x_antes == x and y_antes == y):
        cont = 0
        if mapa_original[x-1][y] == b'0':
            cont = cont+1
        if mapa_original[x+1][y] == b'0': 
            cont = cont+1
        if mapa_original[x][y-1] == b'0':
            cont = cont+1
        if mapa_original[x][y+1] == b'0':
            cont = cont+1

        rnd = random.random()
        if cont <= 2 and falhou == False:
            rnd = p[2]
        
        #cima
        if rnd < 0.25:
            if mapa_original[x-1][y] == b'0':
                x = x - 1
        #baixo
        elif rnd < 0.5:
            if mapa_original[x+1][y] == b'0': 
                x = x + 1
        #esquerda
        elif rnd < 0.75:
            if mapa_original[x][y-1] == b'0':
                y = y - 1
        #direita
        else:
            if mapa_original[x][y+1] == b'0': 
                y = y + 1
        falhou = True
    return x, y, rnd