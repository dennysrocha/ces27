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
for i in range(0, n_linhas):
    for j in range(0, n_colunas):
        if mapa_original[i][j] != b'0':
            mapaint[i][j] = -1

# Acende um bloco de "n_em_volta" casas em todas as direções da posicao
def acende_luzes_em_volta(p, cont):
    global mapa
    if cont >= 1:
        if mapa[p[0] - 1][p[1]] == b'0':
            mapa[p[0] - 1][p[1]] = b'+'
            acende_luzes_em_volta((p[0] - 1, p[1]), cont - 1)
        if mapa[p[0] + 1][p[1]] == b'0':
            mapa[p[0] + 1][p[1]] = b'+'
            acende_luzes_em_volta((p[0] + 1, p[1]), cont - 1)
        if mapa[p[0]][p[1] - 1] == b'0': 
            mapa[p[0]][p[1] - 1] = b'+' 
            acende_luzes_em_volta((p[0], p[1] - 1), cont - 1)
        if mapa[p[0]][p[1] + 1] == b'0':
            mapa[p[0]][p[1] + 1] = b'+'  
            acende_luzes_em_volta((p[0], p[1] + 1), cont - 1)

# Acende o caminho mais percorrido a partir daquela posicao
def acende_luzes_caminho_otimizado(p, cont):
    global mapa
    global mapaint
    sentido = "cima"
    maior = mapaint[p[0]-1][p[1]]

    #baixo
    if maior < mapaint[p[0]+1][p[1]]:
        sentido = "baixo"
        maior = mapaint[p[0]+1][p[1]]

    #esquerda
    if maior < mapaint[p[0]][p[1]-1]:
        sentido = "esquerda"
        maior = mapaint[p[0]][p[1]-1]

    #direita
    if maior < mapaint[p[0]][p[1]+1]:
        sentido = "direita"
        maior = mapaint[p[0]][p[1]+1]

    if cont >= 1:
        if sentido == "cima":
            if mapa[p[0] - 1][p[1]] != b'P':
                mapa[p[0] - 1][p[1]] = b'+'
            acende_luzes_caminho_otimizado((p[0] - 1, p[1]), cont - 1)
        elif sentido == "baixo":
            if mapa[p[0] + 1][p[1]] != b'P':
                mapa[p[0] + 1][p[1]] = b'+'
            acende_luzes_caminho_otimizado((p[0] + 1, p[1]), cont - 1)
        elif sentido ==  "esquerda":
            if mapa[p[0]][p[1] - 1] != b'P': 
                mapa[p[0]][p[1] - 1] = b'+' 
            acende_luzes_caminho_otimizado((p[0], p[1] - 1), cont - 1)
        elif sentido ==  "direita":
            if mapa[p[0]][p[1] + 1] != b'P':
                mapa[p[0]][p[1] + 1] = b'+'  
            acende_luzes_caminho_otimizado((p[0], p[1] + 1), cont - 1)

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
            elif matriz[i][j] == b'+':
                print ("+", end = '')
            else:
                print ("x", end = '')
        print() 

# Atualiza mapaint
def atualiza_mapaint(p):
    global mapaint
    mapaint[p[0], p[1]] = mapaint[p[0], p[1]] + 1 

# Ilumina mapa
def ilumina_mapa():
    global mapa
    for i in range(0, n_linhas):
        for j in range(0, n_colunas):
            if mapa[i][j] == b'P':
                p = (i, j)
                acende_luzes_em_volta(p, 16)
                acende_luzes_caminho_otimizado(p, 24)

# Rendimento em %
rendimento = 0
cont_rendimento = 0
def calcula_rendimento():
    global rendimento
    global cont_rendimento
    global mapa
    total_luzes = 0
    luzes_acesas = 0
    for i in range(0, n_linhas):
        for j in range(0, n_colunas):
            if mapa[i][j] == b'+' or mapa[i][j] == b'P':
                luzes_acesas = luzes_acesas + 1
                total_luzes = total_luzes + 1
            elif mapa[i][j] == b'0':
                total_luzes = total_luzes + 1
    rendimento = (rendimento*cont_rendimento + 100*(total_luzes-luzes_acesas)/total_luzes) / (cont_rendimento + 1)
    cont_rendimento = cont_rendimento + 1
    return rendimento

# teste
aux = 0
p1 = (5,5,random.random())
p2 = (5,106,random.random())
p3 = (34,13,random.random())
p4 = (34,94,random.random())
p5 = (17,61,random.random())
atualiza_mapaint(p1)
atualiza_mapaint(p2)
atualiza_mapaint(p3)
atualiza_mapaint(p4)
atualiza_mapaint(p5)
while aux < 6000:
    system('cls') 
    p1 = movimenta(p1)
    p2 = movimenta(p2)
    p3 = movimenta(p3)
    p4 = movimenta(p4)
    p5 = movimenta(p5)
    atualiza_mapaint(p1)
    atualiza_mapaint(p2)
    atualiza_mapaint(p3)
    atualiza_mapaint(p4)
    atualiza_mapaint(p5)
    mapa[p1[0]][p1[1]] = 'P'
    mapa[p2[0]][p2[1]] = 'P'
    mapa[p3[0]][p3[1]] = 'P'
    mapa[p4[0]][p4[1]] = 'P'
    mapa[p5[0]][p5[1]] = 'P'
    ilumina_mapa()
    print_map(mapa)
    print("Rendimento = ", "%.1f" % calcula_rendimento(), " % de economia")
    reiniciar_mapa()    
    time.sleep(0.05)    
    aux = aux + 1
