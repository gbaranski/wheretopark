<!-- # Introduction

## What problem are we solving?

*Search traffic* is traffic that occurs only when looking for a parking space.

Research in Krakow shows that **the share of *search traffic* in total traffic** ranges from **8-10%**, and in other foreign cities, this value can reach **70%** [2].

93.7% of drivers said that in the case when there is no free space on the parking lot, they **circulate in the hope of its release**, and only the remaining 6.3% look for a place on the nearby parking lot [1].

In Los Angeles alone, it was estimated that every year *search traffic* contributes to the production of **730 tons of carbon dioxide** and the consumption of **170 thousand liters of fuel** [3].

## Our solution

*Where To Park* aims to streamline the parking process, and at the same time reduce the above-mentioned search traffic, by giving drivers easy access to information about the number of free spaces on nearby parking lots, so that instead of circling around the parking lot, they can immediately drive to the one next to it, because in the application they can see that there will be free spaces there.

<a href="preview.webp"><img src="preview.webp"/></a>
<div style="text-align: right;">
  <font size="2">App is available for browsers, iOS and Android</font>
</div>
<br/>

# Where do we get the data?

### Public data

Selected agglomerations offer open access to data from existing parking systems.
One of them is Trójmiasto, which through the Tristar system gives access to almost 20 parking lots.

### CCTV cameras

We use artificial intelligence to determine the number of free parking spaces.
Depending on the parking lot and weather conditions, the detection rate of a free space is about 70-90%, with room for improvement, for which we require more data.

Below are examples of visualizations of the system's functioning.

Program automatically marks <span style="color:green">green</span> free spaces, and <span style="color:red">red</span> occupied.


#### Swimming pool parking lot in Kłodzko
<a href="/visualisation/basen_klodzko-1.webp" target="_blank">
    <img src="/visualisation/basen_klodzko-1.webp" width="400px">
</a>

#### "Królewski" Parking lot in Gdańsk

<a href="/visualisation/krolewski_gdansk-1.webp" target="_blank">
    <img src="/visualisation/krolewski_gdansk-1.webp" width="400px">
</a>


# Possible benefits

Drivers using *Where To Park* will be able to find a free parking space faster, which will reduce the amount of *search traffic* and thus reduce the amount of CO2 emissions and fuel consumption.


Use of CCTV cameras is a small intrusion into the infrastructure. We use existing cameras and no special sensors are required.

## For the city

- New possibilities for the city related to the collected data, including:
  - History of free spaces on certain days and hours, as well as any **statistics**, and **trend analysis** regarding drivers. These information can help in choosing the right place to create new parking spaces.
  - Detection of **illegally parked** vehicles.
- Optimization of road traffic
    - **Increase in the occupancy of available parking space**, drivers could start using less known, or recently opened parking lots.
    - **Decrease in the required movements to find a parking space**, and at the same time **decreased overall traffic**.
- Promotion of the city as **innovative and modern**.

## For private parking lots

- More customers, thanks to the promotion of the parking lot in the application.
- In the future, the possibility of **reserving parking spaces** directly in the application, instead phone calling.

## For drivers

- Time and money savings
    - During peak hours in Krakow, the share of free spaces is only 5%, so drivers **lost an average of 10 minutes looking for a free parking space**, with a maximum of 45 minutes [1].
    - **The cost of operating a car, and the cost of time spent on finding a parking space** is estimated at about 9 PLN, not including the additional cost of parking in a paid parking zone (SPP) [1].
    - A study in Los Angeles says that the annual loss of time is at the level of 95000 hours [4].
- **Disabled people** can easily check which parking lots have specially designated spaces for them, and whether these spaces are free.
- **Predictions of availability**. The user has access to the statistics of the occupancy of places on parking lots at a given time and on a given day. For example, on Friday at 5 pm, the parking lot is usually 80% full, and on Saturday at 2 pm, it is 50% full.


# Mobile Application

## Features

- List of parking lots nearby.
- Map with nearby parking lots.
- Displaying data about a given parking lot.

Parking lot data contains:
- Real-time availability of places
- Total number of parking places.
- Opening hours and pricing.
- Contact info.
- Additional comments, e.g. "student card allows for 50% discount on monthly subscription"

# References
1. [Duda-Wiertel U. Konsekwencje zmiany dostępności przestrzeni parkingowej we wrażliwych obszarach centrów miast, Politechnika Krakowska, Kraków 2018.](https://yadda.icm.edu.pl/baztech/element/bwmeta1.element.baztech-51e98198-f84e-4b3f-8a4a-2492f1fef1df)
2. [Duda-Wiertel U. Search traffic w obszarach z deficytem miejsc postojowych, Politechnika Krakowska, Kraków 2021.](https://yadda.icm.edu.pl/baztech/element/bwmeta1.element.baztech-d8bf4176-2ea0-4571-b5a5-bdc066343c33)
3. [Aliedani, Ali & Loke, Seng & Desai, Aniruddha & Desai, Prajakta. (2016). Investigating Vehicle-to-Vehicle Communication for Cooperative Car Parking: the CoPark Approach.](https://www.researchgate.net/publication/307534127_Investigating_Vehicle-to-Vehicle_Communication_for_Cooperative_Car_Parking_the_CoPark_Approach) -->
