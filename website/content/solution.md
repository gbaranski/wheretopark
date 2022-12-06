---
title: "Nasze rozwiązanie"
weight: 3
draft: true
---

# *Search traffic* - problem z parkowaniem w dużych miastach

*Search traffic* to ruch związany wyłącznie z poszukiwaniem miejsca parkingowego.

93.7% kierowców stwierdziło, że w sytuacji gdy obok miejsca docelowego nie ma wolnego miejsca parkingowego, **krążą w oczekiwaniu na jego zwolnienie**, a tylko pozostałe 6.3% szuka miejsca na pobliskim parkingu [1]. 

W samym Los Angeles oszacowano, że każdego roku search traffic przyczynia się do wytworzenia **730 ton dwutlenku węgla** i zużycia **170 tysięcy litrów paliwa** [4].

# Nasze rozwiazanie

*Where To Park* ułatwia kierowcom znalezienie wolnego miejsca parkingowego w przystępny i intuicijny sposób. 

Aplikacja jest dostępna w przeglądarce, na iOS'ie i Androidzie.

# Skąd zbieramy dane?

### Publiczne dane

Przykładowo trójmiasto, oferuje otwarty dostęp do danych pochodzących z systemu Tristar, czyli na ten moment prawie 20 parkingów.

### Kamery monitoringowe
Korzystamy ze sztucznej inteligencji w celu określenia ilości wolnych miejsc parkingowych. 
Rozwiązanie aktualnie testujemy na trzech prywatnych parkingach w Gdańsku i w Kłodzku. 
Zbadana skuteczność detekcji wynosi na ten moment około 90%.

Poniżej przykładowe wizualizacje z działania systemu:

#### Parking przy basenie w Kłodzku
<a href="/visualisation/basen_klodzko-1.jpeg" target="_blank">
    <img src="/visualisation/basen_klodzko-1.jpeg" width="400px">
</a>

#### Parking Królewski w Gdańsku

<a href="/visualisation/krolewski_gdansk-1.jpeg" target="_blank">
    <img src="/visualisation/krolewski_gdansk-1.jpeg" width="400px">
</a>


# Korzyści z rozwiązania

Kierowcy korzystający z naszej aplikacji chętniej zamiast krążenia wokół parkingu, powodując *search traffic*, wybierają wolne miejsce parkingowe, które znajdą w aplikacji.

Rozwiązanie korzystające z kamer monitoringowych to mała ingerencja w infrastrukturę, nie jest wymagany montaż specjalnych czujników, które też nie wszędzie są możliwe do zamontowania.

## Dla miasta

- Nowe możliwości dla miasta związane z posiadaniem zebranych danych, w tym:
  - Historia wolnych miejsc w danych dniach i godzinach, a także wszelkie **statystyki**, i **analizy trendów** dot. kierowców. Pomóc to może w bardziej trafnym wyborze miejsca do utworzenia nowych miejsc parkingowych.
  - Wykrywanie **nieprawidłowo zaparkowanych** pojazdów.
- Optymalizacja ruchu drogowego
  - **Zwiększenie zajętości dostępnej przestrzeni parkingowej**, kierowcy mogliby zacząć korzystać z mniej znanych, lub dopiero co otworzonych parkingów.
  - **Zmniejszenie wymaganych przemieszczeń w celu znalezienia miejsca parkingowego**, a przy tym **zmniejszone ogólne natężenie ruchu**, w pracy [2] autor mówi o wpływie nawet na **8-10% ruchu ogólnego** w przypadku Krakowa, a w przypadku innych zagranicznym miast, ta wartość sięgała **nawet 70%**.
- Promocja miasta jako **innowacyjne i nowoczesne**.

## Dla kierowców
- Oszczędność czasu i pieniędzy
    - Kierowcy w Krakowie **tracili średnio około 10 minut na znalezienie dostępnego miejsca parkingowego**, przy czym maksymalny czas wynosił około 45 minut [1].
    - **Koszt eksploatacyjny samochodu, i koszt czasu poświęcony na jednorazowe szukanie miejsca parkingowego** wynosi szacunkowo około 9zł, nie wliczając w to dodatkowego kosztu samego parkowania w SPP [1].
    - Badanie w Los Angeles mówi o rocznej stracie czasu na poziomie 95000 godzin [4].
- **Osoby niepełnosprawne**, mogą prosto w aplikacji sprawdzić które parkingi maja dla nich specjalnie wyznaczone miejsca, i czy są te miejsca wolne. 
- **Predykcje dostępności** - widoczna dla użytkownika statystyka mówiącą o tym, że przykładowo w piątek o godzinie 10 zazwyczaj ciężko znaleźć wolne miejsce na danym parkingu, a dzięki tej informacji mozna zaplanowac trasę dzień wczesniej.


W przyszłości planujemy także dodatkowe funkcjonalności takie jak np. **rezerwacja miejsc parkingowych**, czy **płatności za postój** w aplikacji.


# Przypisy
1. [Duda-Wiertel U. Konsekwencje zmiany dostępności przestrzeni parkingowej we wrażliwych obszarach centrów miast, Politechnika Krakowska, Kraków 2018.](https://yadda.icm.edu.pl/baztech/element/bwmeta1.element.baztech-51e98198-f84e-4b3f-8a4a-2492f1fef1df)
2. [Duda-Wiertel U. Search traffic w obszarach z deficytem miejsc postojowych, Politechnika Krakowska, Kraków 2021.](https://yadda.icm.edu.pl/baztech/element/bwmeta1.element.baztech-d8bf4176-2ea0-4571-b5a5-bdc066343c33)
3. [Rehman, Mujeeb & Shah, Munam. (2017). A smart parking system to minimize searching time, fuel consumption and CO2 emission. 1-6. 10.23919/IConAC.2017.8082088.](https://www.researchgate.net/publication/320826298_A_smart_parking_system_to_minimize_searching_time_fuel_consumption_and_CO2_emission)
4. [Aliedani, Ali & Loke, Seng & Desai, Aniruddha & Desai, Prajakta. (2016). Investigating Vehicle-to-Vehicle Communication for Cooperative Car Parking: the CoPark Approach.](https://www.researchgate.net/publication/307534127_Investigating_Vehicle-to-Vehicle_Communication_for_Cooperative_Car_Parking_the_CoPark_Approach)
5. [Richard Arnott & Tilmann Rave & Ronnie Schöb, 2005. "Alleviating Urban Traffic Congestion," MIT Press Books, The MIT Press, edition 1, volume 1, number 0262012197.](https://ideas.repec.org/b/mtp/titles/0262012197.html)