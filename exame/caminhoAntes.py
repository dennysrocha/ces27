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

# Matriz de inteiros que indica posicoes mais percorridas
mapaint = np.zeros(dim)
for i in range(0, n_linhas - 1):
    for j in range(0, n_colunas - 1):
        if mapa[i][j] == "0":
            mapaint[i][j] = 0
        else:
            mapaint[i][j] = -1

# Acende um bloco de "n_em_volta" casas em todas as direções da posicao
def acende_luzes_em_volta(x , y):
    n_em_volta = 4
    for i in range(x - n_em_volta, x + n_em_volta):
        for j in range(y - n_em_volta, y + n_em_volta):
            #aceso
            if mapa[i][j] == "0":
                mapa[i][j] = "1"

# Acende o caminho mais percorrido a partir daquela posicao
def acende_luzes_caminho_otimizado(x , y, cont):
    mapa[x][y] = "1"
    sentido = "cima"
    maior = mapaint[x-1][y]

    #baixo
    if maior < mapaint[x+1][y]:
        sentido = "baixo"
        maior = mapaint[x+1][y]

    #esquerda
    if maior < mapaint[x][y-1]:
        sentido = "esquerda"
        maior = mapaint[x][y-1]

    #direita
    if maior < mapaint[x][y+1]:
        sentido = "direita"
        maior = mapaint[x][y+1]

    if cont > 1:
        if sentido == "cima":
            acende_luzes_caminho_otimizado(x - 1, y, cont - 1)
        elif sentido == "baixo":  
            acende_luzes_caminho_otimizado(x + 1, y, cont - 1)
        elif sentido ==  "esquerda":  
            acende_luzes_caminho_otimizado(x, y - 1, cont - 1)
        elif sentido ==  "direita":  
            acende_luzes_caminho_otimizado(x, y + 1, cont - 1)

# Reinicia o mapa
def reiniciar_mapa():
    global mapa
    mapa = deepcopy(mapa_original)

# Movimento aleatório do usuário
def movimenta(p):
    x = p[0]
    y = p[1]
    x_antes = x
    y_antes = y

    # Escolhe algum zero inicial
    """if mapa_original[x-1][y] == b'0':
        x = x - 1
    elif mapa_original[x+1][y] == b'0': 
        x = x + 1
    elif mapa_original[x][y+1] == b'0': 
        y = y + 1
    elif mapa_original[x][y-1] == b'0':
        y = y - 1"""
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
            #print("entrei")
            #time.sleep(1)
        
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

# Print do mapa
def print_map(matriz):
    for i in range (0, n_linhas):  
        for j in range (0, n_colunas):
            if matriz[i][j] == b'.':
                print (".", end = '')
            elif matriz[i][j] == b'0':
                print (" ", end = '')
            elif matriz[i][j] == b'P':
                print ("☺", end = '')
            else:
                print ("x", end = '')
        print() 


# teste
aux = 0
p1 = (5,5,random.random())
p2 = (5,106,random.random())
p3 = (34,13,random.random())
p4 = (34,94,random.random())
p5 = (17,61,random.random())
while aux < 6000:
    system('cls') 
    p1 = movimenta(p1)
    p2 = movimenta(p2)
    p3 = movimenta(p3)
    p4 = movimenta(p4)
    p5 = movimenta(p5)
    #print("x = ", x)
    #print("y = ", y)
    mapa[p1[0]][p1[1]] = 'P'
    mapa[p2[0]][p2[1]] = 'P'
    mapa[p3[0]][p3[1]] = 'P'
    mapa[p4[0]][p4[1]] = 'P'
    mapa[p5[0]][p5[1]] = 'P'
    print_map(mapa)
    #print_map(mapa_original)
    reiniciar_mapa()
    #print(aux)
    time.sleep(0.05)    
    aux = aux + 1
