+++
title = "Where To Park"
description = "Where To Park to aplikacja pomagająca kierowcom odnaleźć wolne miejsce parkingowe, wykorzystująca techniki sztucznej inteligencji i analizy danych w czasie rzeczywistym."
sort_by = "weight"

[extra]
app_store_badge = "https://tools.applemediaservices.com/api/badges/download-on-the-app-store/white/pl-pl?size=250x83&amp;releaseDate=1668988800&h=ef8147fbbaa63160018d2a5be2e2565e"
+++


# Wprowadzenie

## Jaki problem rozwiązujemy?

*Search traffic* to ruch wynikający wyłącznie z poszukiwania miejsca parkingowego.

Z badań w Krakowie wynika, że **udział *search traffic* w ruchu ogólnym** oscyluje w okolicach **8-10%**, a w przypadku innych zagranicznych miast, ta wartość może sięgać nawet **70%** [2].

93.7% kierowców stwierdziło, że w sytuacji gdy na parkingu nie ma wolnego miejsca, **krążą w oczekiwaniu na jego zwolnienie**, a tylko pozostałe 6.3% szuka miejsca na pobliskim parkingu [1]. 


W samym Los Angeles oszacowano, że każdego roku *search traffic* przyczynia się do wytworzenia **730 ton dwutlenku węgla** i zużycia **170 tysięcy litrów paliwa** [3].

## Nasze rozwiazanie

*Where To Park* ma na celu usprawnić proces parkowania, i przy tym obniżyć wyżej wymieniony search traffic, poprzez danie kierowcom łatwego dostępu do informacji o ilości wolnych miejsc na pobliskich parkingach, przez co zamiast krążenia wokół parkingu, mogą oni od razu pojechać ten obok, ponieważ w aplikacji mogą zobaczyć ze znajdą się na nim wolne miejsca.

<img src="/pl/preview.jpg"/>
<div style="text-align: right;">
  <font size="2">Aplikacja jest dostępna na przeglądarki, iOS i Android.</font>
</div>
<br/>

# Skąd zbieramy dane?

### Publiczne dane

Wybrane aglomeracje oferują otwarty dostęp do danych pochodzących z istniejących systemów parkingowych. 
Jednym z nich jest Trójmiasto, które za pośrednictwem systemu Tristar daje dostęp do prawie 20 parkingów.

### Kamery monitoringowe
Korzystamy ze sztucznej inteligencji w celu określenia liczby wolnych miejsc parkingowych. 
Rozwiązanie aktualnie testujemy na trzech prywatnych parkingach w Gdańsku oraz Kłodzku. 
W zależności od parkingu oraz warunków atmosferycznych, skuteczność detekcji wolnego miejsca wynosi około 70-90%, z miejscem na poprawę, do której wymagamy większej ilości danych.

Niżej przykładowe wizualizacje funkcjonowania systemu.

Program automatycznie oznaczył na <span style="color:green">zielono</span> miejsca wolne, a na <span style="color:red">czerwono</span> zajęte.

#### Parking przy basenie w Kłodzku
<a href="/visualisation/basen_klodzko-1.webp" target="_blank">
    <img src="/visualisation/basen_klodzko-1.webp" width="400px">
</a>

#### Parking Królewski w Gdańsku

<a href="/visualisation/krolewski_gdansk-1.webp" target="_blank">
    <img src="/visualisation/krolewski_gdansk-1.webp" width="400px">
</a>


# Korzyści z rozwiązania

Kierowcy korzystający z naszej aplikacji, zamiast krążenia wokół parkingu (powodując *search traffic*), wybierają wolne miejsce parkingowe sugerowane przez naszą aplikację.


Rozwiązanie korzystające z kamer monitoringowych to mała ingerencja w infrastrukturę. Korzystamy z istniejących kamer i nie jest wymagany montaż specjalnych czujników.

## Dla miasta

- Nowe możliwości dla miasta związane z posiadaniem zebranych danych, w tym:
  - Historia wolnych miejsc w danych dniach i godzinach, a także wszelkie **statystyki**, i **analizy trendów** dot. kierowców. Informacje te mogą pomóc w trafnym wyborze miejsca do utworzenia nowych miejsc parkingowych.
  - Wykrywanie **nieprawidłowo zaparkowanych** pojazdów.
- Optymalizacja ruchu drogowego
  - **Zwiększenie zajętości dostępnej przestrzeni parkingowej**, kierowcy mogliby zacząć korzystać z mniej znanych, lub dopiero co otworzonych parkingów.
  - **Zmniejszenie wymaganych przemieszczeń w celu znalezienia miejsca parkingowego**, a przy tym **zmniejszone ogólne natężenie ruchu**.
- Promocja miasta jako **innowacyjne i nowoczesne**.

## Dla prywatnych parkingów
- Więcej klientów, ze względu na promocje parkingu w aplikacji.
- W przyszłosci, możliwość rezerwacji miejsc parkingowych bezpośrednio w aplikacji, zamiast przez telefon. 

## Dla kierowców
- Oszczędność czasu i pieniędzy
    - W godzinach szczytu w Krakowie udział wolnych miejsc wynosi tylko 5%, przez co kierowcy **tracili średnio 10 minut na znalezienie dostępnego miejsca parkingowego**, przy czym maksymalny czas wynosił około 45 minut [1].
    - **Koszt eksploatacyjny samochodu, i koszt czasu poświęcony na jednorazowe szukanie miejsca parkingowego** wynosi szacunkowo około 9zł, nie wliczając w to dodatkowego kosztu samego parkowania w strefie płanego parkowania (SPP) [1].
    - Badanie w Los Angeles mówi o rocznej stracie czasu na poziomie 95000 godzin [4].
- **Osoby niepełnosprawne**, mogą w prosty sposób sprawdzić, które parkingi maja dla nich specjalnie wyznaczone miejsca, i czy te miejsca są wolne. 
- **Predykcje dostępności**. Użytkownik ma dostęp do statystyki zajętości miejsc na parkinach o danej godzine i danego dnia. Przykładowo w piątek o godzinie 10 zazwyczaj ciężko znaleźć wolne miejsce na danym parkingu, a dzięki tej informacji mozna zaplanowac trasę dzień wczesniej.


W przyszłości planujemy także dodatkowe funkcjonalności takie jak np. **rezerwacja miejsc parkingowych**, czy **płatności za postój** w aplikacji.

# Aplikacja

## Funkcje
- Lista parkingów w pobliżu.
- Mapa z pobliskimi parkingami.
- Wyświetlanie danych na temat danego parkingu.

Dane na temat parkingu zawierają:
- Dostępność miejsc w czasie rzeczywistym
- Łączna ilośc miejsc
- Cennik oraz godziny otwarcia.
- Dane kontaktowe
- Dodatkowe komentarze, np. "legitymacja studencka pozwala na 50% zniżki przy subskrypcji miesięcznej"

## Zrzuty ekranu

<div class="demonstration-media">
    <img class="demonstration-image" src="screenshots/main.webp" alt="screenshot-main">
    <img class="demonstration-image" src="screenshots/selected.webp" alt="screenshot-selected">
    <img class="demonstration-image" src="screenshots/expanded.webp" alt="screenshot-expanded">
    <img class="demonstration-image" src="screenshots/expanded-2.webp" alt="screenshot-expanded-2">
</div>


# Przypisy
1. [Duda-Wiertel U. Konsekwencje zmiany dostępności przestrzeni parkingowej we wrażliwych obszarach centrów miast, Politechnika Krakowska, Kraków 2018.](https://yadda.icm.edu.pl/baztech/element/bwmeta1.element.baztech-51e98198-f84e-4b3f-8a4a-2492f1fef1df)
2. [Duda-Wiertel U. Search traffic w obszarach z deficytem miejsc postojowych, Politechnika Krakowska, Kraków 2021.](https://yadda.icm.edu.pl/baztech/element/bwmeta1.element.baztech-d8bf4176-2ea0-4571-b5a5-bdc066343c33)
3. [Aliedani, Ali & Loke, Seng & Desai, Aniruddha & Desai, Prajakta. (2016). Investigating Vehicle-to-Vehicle Communication for Cooperative Car Parking: the CoPark Approach.](https://www.researchgate.net/publication/307534127_Investigating_Vehicle-to-Vehicle_Communication_for_Cooperative_Car_Parking_the_CoPark_Approach)

