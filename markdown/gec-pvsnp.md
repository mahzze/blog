Title: P vs NP: Complexidade computacional.
Parent: GEC - Guia dos entusiastas da ciência.
Order: 1
Slug: gec-pvsnp
MetaPropertyTitle: P vs NP
MetaDescription: Complexidade computacional e o problema mais importante do milênio
MetaOgURL: https://www.mahzze.dev/
---

O seguinte texto é um artigo que escrevi para o guia dos entusiastas da ciência (GEC). Este texto especificamente ainda não foi (e talvez nem seja) publicado. Outros como ele podem ser encontrados no blog: https://gec.proec.ufabc.edu.br

---
## Já ouviu falar nos sete problemas do milênio?

O título chamou atenção e você está aí pensando em quais, dentre tantos possíveis, seriam os sete problemas do milênio? De fato, adivinhar os principais problemas do futuro da humanidade não é uma tarefa fácil. Por isso, este artigo marca o começo de uma série de artigos que dissertam sobre esses problemas. Nesse texto, você vai conhecer um dentre os sete: o problema P vs. NP. Ainda sem solução, ele é considerado o mais importante da ciência da computação. Ele une as mais diversas áreas da matemática, tanto pura como aplicada, abrangendo áreas como a teoria de grafos, probabilidade, topologia, álgebra e teoria dos números, revelando conexões não-triviais entre elas. 

Formalizado em 1971 e advindo do campo da complexidade computacional, que é o ramo da ciência da computação que estuda a eficiência de algoritmos, P vs NP é um problema que é relativamente simples de entender conceitualmente. Mas antes de entender o problema, é necessário saber o que é um algoritmo e o que são as classes P e NP.

Vamos começar por algoritmos: em termos simples, um algoritmo é um conjunto de instruções que são executadas uma após a outra, como uma receita culinária. Em computadores, essas instruções são codificadas em 0’s e 1’s, isto é, são codificadas no sistema binário, para eventualmente serem executadas pelo computador. Para facilitar a escrita de códigos, foram criadas as linguagens de programação, que são ferramentas que permitem que um algoritmo seja escrito em forma de código em uma linguagem simples, para então ser traduzido para binário. Um exemplo de algoritmo bem simples que será utilizado ao longo do texto é o seguinte código em Python de função exemplo que verifica se os números de uma lista são pares ou ímpares:

<center>
<div class="quote-box">
```python
def verifica_par(lista_valores):
    paridades = []
    for valor in lista_valores:
        if valor % 2 == 0:
            paridades.append(True)
        else:
            paridades.append(False)
    return paridades
def main():
    linha_teste = [0, 1, 2, 3, 4, 5, 6, 7]
    print(verifica_par(linha_teste))
```
</div>
</center>

Caso não saiba programar, não se assuste, o código é simples de entender. É criada uma função – na imagem: ***verifica_par(lista_valores)*** -, que recebe uma lista de valores numéricos e verifica se eles são pares, e anota o valor da paridade de cada valor em uma outra lista, anotando **True (Verdadeiro)** caso o número seja par, e **False (Falso)** caso o número seja ímpar.

Ao executarmos a função com os números inteiros de 0 a 7, a lista de paridade que é criada pela função é composta por valores booleanos (True e False) que alternam entre si, como pode ser visto no vetor abaixo.

<center>
<div class="quote-box">
```
[True, False, True, False, True, False, True, False]
```
</div>
</center>

Vamos agora expor um pouco da natureza binária do `nosso` código. Mais especificamente, vamos verificar os valores numéricos de True e False. Isso é feito no código abaixo:

<center>
<div class="quote-box">
```
print("O valor numérico de True é: ", int(True))
print("O valor numérico de False é: ", int(False))
```
</div>
</center>

Executar este trecho de código resulta em:

<center>
<div class="quote-box">
```
O valor numérico de True é: 1
O valor numérico de False é: 0
```
</div>
</center>

Isso são valores binários! Isso significa que, por debaixo dos panos, a lista de paridades criada pela função verifica_par é apenas uma lista composta por itens com valores de 0 e 1. Com isso em mente, do ponto de vista matemático, o nosso exemplo pode ser escrito como:

$ verifica_par([0,1,2,3,4,5,6,7]) \xrightarrow[] [1,0,1,0,1,0,1,0] $

É importante lembrar que os valores numéricos da lista de valores que são passados para a função *verifica_par* também podem ser descritos de forma binária, que é como o computador os processa. Falando da função *verifica_par*, ela pode ser denotada da seguinte forma:
