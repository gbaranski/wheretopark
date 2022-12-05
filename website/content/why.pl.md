---
title: "Dlaczego my?"
description: "Dlaczego my?"
weight: 3
draft: true
---

# Problem z tzw. “search traffic”
*Search traffic* to ruch związany wyłącznie z poszukiwaniem miejsca parkingowego.

93.7% kierowców stwierdziło ze w sytuacji gdy u celi podróży nie ma wolnego miejsca parkingowego, **krążą w oczekiwaniu na zwolnienie miejsca**, a tylko pozostałe 6.7% szuka na pobliskim parkingu [1]. 

Skutki ekologiczne szukania wolnego miejsca parkingowego w samym Los Angeles oszacowano w na **730 ton dwutlenku węgla** i **170 tysięcy litrów paliwa** rocznie [4].

# Nasze rozwiazanie

Aplikacja *Where To Park* została stworzona z myślą o kierowcach, którzy aktywnie korzystają z transportu drogowego.
Ma zadanie ułatwić znalezienie wolnego miejsca parkingowego, w przystępny i intuicyjny sposób.

Aplikacja juz teraz jest dostępna w przeglądarce, na iOS'ie i Androidzie.

# Skąd zbieramy dane?

### Publiczne dane

W Trójmieście jest to np system Tristar, który publicznie udostępnia dane o wolnych miejscach na prawie 20 parkingach.


### Kamery monitoringowe
Korzystamy ze sztucznej inteligencji w celu określenia czy dane miejsce parkingowe jest wolne. 
Rozwiązanie aktualnie testujemy na trzech prywatnych parkingach w Gdańsku i w Kłodzku. 

Tutaj przykładowe wizualizacje z działania systemu:

#### Parking przy basenie w Kłodzku
<a href="/visualisation/basen_klodzko-1.jpeg" target="_blank">
    <img src="/visualisation/basen_klodzko-1.jpeg" width="400px">
</a>

#### Parking Królewski w Gdańsku

<a href="/visualisation/krolewski_gdansk-1.jpeg" target="_blank">
    <img src="/visualisation/krolewski_gdansk-1.jpeg" width="400px">
</a>


# Zalety rozwiązania

Aplikacja może zachęcić kierowców aby pojechali na pobliski parking, ponieważ mają pewność ze na na tym parkingu znajdą się dla nich wolne miejsca, a za tym idą dodatkowe korzyści takie jak mniejsze natężenie ruchu i mniej zbędnych kilometrów, czyli też i mniej zanieczyszczeń w powietrzu.

Rozwiązanie korzystające z kamer monitoringowych to mała ingerencja w infrastrukturę, nie jest wymagany montaż specjalnych czujników, które też nie wszędzie są możliwe do zamontowania.

## Dla miasta

- Nowe możliwości dla miasta związane z posiadaniem zebranych danych, w tym:
  - **Analiza trendów** wynikających z ilości wolnych miejsc w danych godzinach, czy dniach. Pomóc to może w bardziej trafnym wyborze miejsca do utworzenia nowych miejsc parkingowych.
  - Wykrywanie **nieprawidłowo zaparkowanych** pojazdów.
- Optymalizacja ruchu drogowego
  - **Zwiększenie zajętości dostępnej przestrzeni parkingowej**, kierowcy mogliby zacząć korzystać z mniej znanych, lub dopiero co otworzonych parkingów.
  - Zmniejszenie wymaganych przemieszczeń w celu znalezienia miejsca parkingowego, a przy tym **zmniejszone ogólne natężenie ruchu**, w pracy [2] autor mówi o wpływie nawet na **8-10% ruchu ogólnego** w przypadku Krakowa, a w przypadku innych zagranicznym miast, ta wartość sięgała **nawet 70%**.
- Promocja miasta jako **innowacyjne i nowoczesne**.
- **Specjalne oznaczenia** dla np. parkingów Park & Ride, które mogą pomoc w dotarciu do użytkowników.

## Dla kierowców
- Oszczędność czasu i pieniędzy
    - Kierowcy w Krakowie **tracili średnio około 10 minut na znalezienie dostępnego miejsca parkingowego**, przy czym maksymalny czas wynosił około 45 minut [1]. Taka aplikacja może drastycznie skrócić ten czas.
    - **Koszt eksploatacyjny samochodu**, i koszt czasu poświęcony na szukanie miejsca parkingowego wynosi szacunkowo okolo 9zł, nie wliczając w to dodatkowego kosztu samego parkowania w SPP. A mnożąc tą wartość przez 20 dni roboczych miesiąca, otrzymujemy wartość **180 zł miesięcznie** [1].
    - Badanie w Los Angeles mówi o rocznej stracie czasu na poziomie 95000 godzin [4].
- **Osoby niepełnosprawne**, mogą prosto w aplikacji sprawdzić które parkingi maja dla nich specjalnie wyznaczone miejsca, i czy są te miejsca wolne. 
- **Predykcje dostępności** bazowane na historii ilości wolnych miejsc, np. widoczna dla użytkownika statystyka mówiącą o tym, że w piątek o godzinie 10 zazwyczaj nie ma ani jednego wolnego miejsca na danym parkingu, przez co może on inaczej zaplanować wcześniej trasę.


W przyszłości planujemy także dodatkowe funkcjonalności takie jak np. **rezerwacja miejsc parkingowych**, czy **płatności za postój** w aplikacji.


# Przypisy
1. [Duda-Wiertel U. Konsekwencje zmiany dostępności przestrzeni parkingowej we wrażliwych obszarach centrów miast, Politechnika Krakowska, Kraków 2018.](https://yadda.icm.edu.pl/baztech/element/bwmeta1.element.baztech-51e98198-f84e-4b3f-8a4a-2492f1fef1df)
2. [Duda-Wiertel U. Search traffic w obszarach z deficytem miejsc postojowych, Politechnika Krakowska, Kraków 2021.](https://yadda.icm.edu.pl/baztech/element/bwmeta1.element.baztech-d8bf4176-2ea0-4571-b5a5-bdc066343c33)
3. [Rehman, Mujeeb & Shah, Munam. (2017). A smart parking system to minimize searching time, fuel consumption and CO2 emission. 1-6. 10.23919/IConAC.2017.8082088.](https://www.researchgate.net/publication/320826298_A_smart_parking_system_to_minimize_searching_time_fuel_consumption_and_CO2_emission)
4. [Aliedani, Ali & Loke, Seng & Desai, Aniruddha & Desai, Prajakta. (2016). Investigating Vehicle-to-Vehicle Communication for Cooperative Car Parking: the CoPark Approach.](https://www.researchgate.net/publication/307534127_Investigating_Vehicle-to-Vehicle_Communication_for_Cooperative_Car_Parking_the_CoPark_Approach)
5. [Richard Arnott & Tilmann Rave & Ronnie Schöb, 2005. "Alleviating Urban Traffic Congestion," MIT Press Books, The MIT Press, edition 1, volume 1, number 0262012197.](https://ideas.repec.org/b/mtp/titles/0262012197.html)