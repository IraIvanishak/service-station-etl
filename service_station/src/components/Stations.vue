<script>
import { API_URL } from '@/api';
import StationDetail from './StationDetail.vue';
import getUnicodeFlagIcon from 'country-flag-icons/unicode'
import { countries } from 'country-flag-icons'
export default {
    data() {
        return {
            stations: [],
            newStation: {
                Address: '',
                CityId: 0
            },
            countryIdMap: {
                "Україна": "UA",
                "США": "US",
                "Канада": "CA",
                "Велика Британія": "GB",
                "Франція": "FR",
                "Німеччина": "DE",
                "Італія": "IT",
                "Іспанія": "ES",
                "Польща": "PL",
                "Австралія": "AU",
                "Нідерланди": "NL",
                "Бельгія": "BE",
                "Швеція": "SE",
                "Швейцарія": "CH",
                "Австрія": "AT",
                "Норвегія": "NO",
                "Данія": "DK",
                "Фінляндія": "FI",
                "Португалія": "PT",
                "Ірландія": "IE",
                "Нова Зеландія": "NZ",
                "Греція": "GR",
                "Японія": "JP",
                "Індія": "IN",
                "Китай": "CN",
                "Бразилія": "BR",
                "Мексика": "MX",
                "Індонезія": "ID",
                "Туреччина": "TR",
                "Шри-Ланка": "LK",
                "Іран": "IR",
                "Таїланд": "TH",
                "Аргентина": "AR",
                "Пакистан": "PK",
                "Бангладеш": "BD",
                "В'єтнам": "VN",
                "Єгипет": "EG",
                "Уганда": "UG",
                "Нігерія": "NG",
                "Кенія": "KE",
                "Гана": "GH",
                "Танзанія": "TZ",
                "Ефіопія": "ET",
                "Саудівська Аравія": "SA",
                "Йорданія": "JO",
                "Марокко": "MA",
                "Алжир": "DZ",
                "Південна Корея": "KR",
                "Грузія": "GE"
            },
            cities: []
        };
    },
    mounted() {
        console.log('Component mounted.');
        this.fetchStations();
        this.listCities();
    },
    methods: {
        goToStationDetail(stationId) {
            console.log(stationId)
            this.$router.push("/stations/" + stationId);
        },
        startEditing(station) {
            station.editing = true;
        },

        cancelEditing(station) {
            station.editing = false;
        },

        async fetchStations() {
            try {
                console.log(countries)
                const response = await fetch(`${API_URL}/stations`);
                const data = await response.json();
                this.stations = data;
                this.stations = data.map(station => {
                    console.log(this.countryIdMap[station.Name_2])
                    return {
                        ...station,
                        flagIcon: getUnicodeFlagIcon(this.countryIdMap[station.Name_2])
                    };
                });
            } catch (error) {
                console.error(error);
            }
        },
        // add new station post http://127.0.0.1:3000/stations
        async addStation() {
            try {
                const response = await fetch(`${API_URL}/stations`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        Address: this.newStation.Address,
                        CityId: Number(this.newStation.CityId),
                    })
                });
                this.fetchStations();

            } catch (error) {
                console.error(error);
            }
        },
        async submitUpdate(station) {
            try {
                const response = await fetch(`${API_URL}/stations`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        Id: station.StationID,
                        Address: station.Address,
                        CityId: station.City
                    })
                });
                if (response.ok) {
                    station.editing = false;
                    this.fetchStations();
                }

            } catch (error) {
                console.error(error);
            }
        },
        async deleteStation(stationId) {
            console.log(stationId)
            try {
                let id = stationId
                const response = await fetch(`${API_URL}/stations/${id}`, {
                    method: 'DELETE',
                });
                this.fetchStations();

                if (response.ok) {
                    // Remove the station from the list
                    this.stations = this.stations.filter(station => station.Id !== stationId);
                }
            } catch (error) {
                console.error(error);
            }
        },
        async listCities() {
            try {
                const response = await fetch(`${API_URL}/city`);
                const data = await response.json();
                this.cities = data;
            } catch (error) {
                console.error(error);
            }
        }
    }
};

</script>

<template>
    <div class="container">
        <h1>Service Stations</h1>
        <div class="grid">
            <div v-for="station in stations" :key="station.StationID" class="station"
                @click="goToStationDetail(station.StationID)">
                <div v-if="!station.editing">
                    <h4>{{ station.Name }} <span v-html="station.flagIcon"></span></h4>
                    <h2><i class="fas fa-map-marker-alt"></i> {{ station.Address }}</h2>
                    <button @click.stop="startEditing(station)" class="button edit-button">
                        <i class="fas fa-edit"></i> Edit
                    </button>
                    <button @click.stop="deleteStation(station.StationID)" class="button delete-button">
                        <i class="fas fa-trash-alt"></i> Delete
                    </button>
                </div>
                <div v-else>
                    <div class="form-group">
                        <label for="edit-address-{{station.StationID}}">Address</label>
                        <input type="text" :id="'edit-address-' + station.StationID" v-model="station.Address"
                            @click.stop>
                    </div>
                    <div class="form-group">
                        <label :for="'edit-city-' + station.StationID">City</label>
                        <select :id="'edit-city-' + station.StationID" v-model="station.City" @click.stop>
                            <option v-for="city in cities" :value="city.CityID" :key="city.CityID">{{ city.Name }}
                            </option>
                        </select>
                    </div>
                    <button @click.stop="submitUpdate(station)">Submit</button>
                    <button @click.stop="cancelEditing(station)">Cancel</button>
                </div>
            </div>
            <div class="station-add">
                <h3>Add a Station</h3>
                <div class="form-group">
                    <label for="new-address">Address</label>
                    <input type="text" id="new-address" v-model="newStation.Address" required @click.stop>
                </div>
                <div class="form-group">
                    <label for="new-city">City</label>
                    <select id="new-city" v-model="newStation.CityId" required @click.stop>
                        <option v-for="city in cities" :value="city.CityID" :key="city.CityID">{{ city.Name }}</option>
                    </select>
                </div>
                <button type="button" @click="addStation">Add Station</button>
            </div>
        </div>
    </div>
</template>

<style scoped>
h1 {
    color: #333;
}

.container {
    padding-right: 100px;
}

.form-group label {
    display: block;
    margin-bottom: 0.5rem;
    /* Adjust as needed */
}

.form-group input,
.form-group select {
    width: 100%;
    /* This ensures the inputs fill the container */
    font-size: 1rem;
    line-height: 1.5;
    border-radius: 0.25rem;
    border: 1px solid #ced4da;
    /* Adjust color as needed */
}

h2 {
    margin: 0;
}

.grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    grid-gap: 16px;
}

.station-add {
    height: 280px;
    background-color: rgba(255, 239, 240, 0.507);
    padding-left: 20px;
    padding-right: 20px;
    padding: 20px;

    border: solid 1px rgba(150, 13, 20, 0.247);
    border-radius: 5px;
}

.station {
    height: 280px;
    background-color: rgb(255, 239, 240);
    padding: 20px;
    border-radius: 5px;
}

h3 {
    color: rgb(1, 1, 10);
    font-size: 18px;
    margin-bottom: 5px;
}

p {
    font-size: 14px;
    margin-bottom: 3px;
}

button {
    z-index: 20;
    padding: 10px 15px;
    background-color: #911247;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 16px;
    margin-right: 14px;
    margin-top: 10px;
}

label {
    width: 200px;
    padding-right: 20px;
}

button:hover {
    background-color: #0056b3;
}


.station button,
.station-add button {
    font-size: 14px;
}

input[type="text"] {
    height: 28px;
    margin: 0;
    margin-bottom: 10px;
}

select {
    height: 28px;
}
</style>