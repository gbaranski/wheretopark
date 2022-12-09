---
title: "Nasze rozwiązanie"
description: "Opis naszego rozwiązania, prezentacja działania sztucznej inteligencji"
weight: 0
extra:
  hidden: true
---

# *Search traffic* - problem z parkowaniem w dużych miastach

*Search traffic* to ruch drogowy wynikający wyłącznie z poszukiwania miejsca parkingowego.

93.7% kierowców stwierdziło, że w sytuacji gdy obok miejsca docelowego nie ma wolnego miejsca parkingowego, **krążą w oczekiwaniu na jego zwolnienie**, a tylko pozostałe 6.3% szuka miejsca na pobliskim parkingu [1]. 

W samym Los Angeles oszacowano, że każdego roku *search traffic* przyczynia się do wytworzenia **730 ton dwutlenku węgla** i zużycia **170 tysięcy litrów paliwa** [4].

# Nasze rozwiazanie

*Where To Park* ułatwia kierowcom znalezienie wolnego miejsca parkingowego w obrębie miejsca docelowego. Dostarcza ona mapę parkingów, liczbę wolnych miejsc, cennik, oraz kontakt. Dane zbierane są z publicznie dostępnych baz danych oraz rozpoznawania obrazu z publicznych kamer. Aplikacja jest dostępna na przeglądarki, iOS i Android.

# Skąd zbieramy dane?

### Publiczne dane

Wybrane aglomeracje oferują otwarty dostęp do danych pochodzących z istniejących systemów parkingowych. 
Jednym z nich jest Trójmiasto, które za pośrednictwem systemu Tristar daje dostęp do prawie 20 parkingów.

### Kamery monitoringowe
Korzystamy ze sztucznej inteligencji w celu określenia liczby wolnych miejsc parkingowych. 
Rozwiązanie aktualnie testujemy na trzech prywatnych parkingach w Gdańsku oraz Kłodzku. 
W zależności od parkingu oraz warunków atmosferycznych, skuteczność detekcji wolnego miejsca wynosi około 70-90%.

Poniżej przykładowe wizualizacje z działania systemu:

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
  - **Zmniejszenie wymaganych przemieszczeń w celu znalezienia miejsca parkingowego**, a przy tym **zmniejszone ogólne natężenie ruchu**, w pracy [2] autor mówi o wpływie nawet na **8-10% ruchu ogólnego** w przypadku Krakowa, a w przypadku innych zagranicznym miast, ta wartość sięgała **nawet 70%**.
- Promocja miasta jako **innowacyjne i nowoczesne**.

## Dla kierowców
- Oszczędność czasu i pieniędzy
    - Kierowcy w Krakowie **tracili średnio około 10 minut na znalezienie dostępnego miejsca parkingowego**, przy czym maksymalny czas wynosił około 45 minut [1].
    - **Koszt eksploatacyjny samochodu, i koszt czasu poświęcony na jednorazowe szukanie miejsca parkingowego** wynosi szacunkowo około 9zł, nie wliczając w to dodatkowego kosztu samego parkowania w strefie płanego parkowania (SPP) [1].
    - Badanie w Los Angeles mówi o rocznej stracie czasu na poziomie 95000 godzin [4].
- **Osoby niepełnosprawne**, mogą w prosty sposób sprawdzić, które parkingi maja dla nich specjalnie wyznaczone miejsca, i czy te miejsca są wolne. 
- **Predykcje dostępności**. Użytkownik ma dostęp do statystyki zajętości miejsc na parkinach o danej godzine i danego dnia. Przykładowo w piątek o godzinie 10 zazwyczaj ciężko znaleźć wolne miejsce na danym parkingu, a dzięki tej informacji mozna zaplanowac trasę dzień wczesniej.


W przyszłości planujemy także dodatkowe funkcjonalności takie jak np. **rezerwacja miejsc parkingowych**, czy **płatności za postój** w aplikacji.


# Przypisy
1. [Duda-Wiertel U. Konsekwencje zmiany dostępności przestrzeni parkingowej we wrażliwych obszarach centrów miast, Politechnika Krakowska, Kraków 2018.](https://yadda.icm.edu.pl/baztech/element/bwmeta1.element.baztech-51e98198-f84e-4b3f-8a4a-2492f1fef1df)
2. [Duda-Wiertel U. Search traffic w obszarach z deficytem miejsc postojowych, Politechnika Krakowska, Kraków 2021.](https://yadda.icm.edu.pl/baztech/element/bwmeta1.element.baztech-d8bf4176-2ea0-4571-b5a5-bdc066343c33)
3. [Rehman, Mujeeb & Shah, Munam. (2017). A smart parking system to minimize searching time, fuel consumption and CO2 emission. 1-6. 10.23919/IConAC.2017.8082088.](https://www.researchgate.net/publication/320826298_A_smart_parking_system_to_minimize_searching_time_fuel_consumption_and_CO2_emission)
4. [Aliedani, Ali & Loke, Seng & Desai, Aniruddha & Desai, Prajakta. (2016). Investigating Vehicle-to-Vehicle Communication for Cooperative Car Parking: the CoPark Approach.](https://www.researchgate.net/publication/307534127_Investigating_Vehicle-to-Vehicle_Communication_for_Cooperative_Car_Parking_the_CoPark_Approach)
5. [Richard Arnott & Tilmann Rave & Ronnie Schöb, 2005. "Alleviating Urban Traffic Congestion," MIT Press Books, The MIT Press, edition 1, volume 1, number 0262012197.](https://ideas.repec.org/b/mtp/titles/0262012197.html)
