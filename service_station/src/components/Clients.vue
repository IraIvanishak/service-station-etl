<template>
    <h2>Clients</h2>

    <div class="client-container">
        <div class="add-client-form">
            <h2>Add New Client</h2>
            <div>
                <input v-model="newClient.FirstName" placeholder="First Name" required>
                <input v-model="newClient.LastName" placeholder="Last Name" required>
                <input v-model="newClient.PhoneNumber" placeholder="Phone Number" required>
                <button @click="addClient" type="submit">Add Client</button>
            </div>
        </div>
        <div v-for="client in clients" :key="client.ClientID" class="client-card">
            <h2><i class="fas fa-user"></i> {{ client.FirstName }} {{ client.LastName }}</h2>
            <p><i class="fas fa-phone"></i> {{ client.PhoneNumber.Valid ? client.PhoneNumber.String : 'No phone number'
                }}</p>
            <div class="car-container">
                <h3><i class="fas fa-car"></i> Cars:</h3>
                <ul class="car-list">
                    <li v-for="car in getCarsForClient(client.ClientID)" :key="car.CarID" class="car-item">
                        <span><i class="fas fa-car-side"></i> {{ car.Model.String }} </span>
                        <span><i class="fas fa-palette"></i> {{ car.Color.String }}</span>
                        <span><i class="fas fa-id-card"></i> {{ car.RegistrationNumber.String }}</span>
                    </li>
                </ul>
            </div>
        </div>
    </div>
</template>
<style scoped>
.add-client-form {
    background-color: #fff1f160;
    border-radius: 10px;
    padding: 20px;
    margin: 10px;
    width: 430px;
    transition: transform 0.3s ease;
    border: solid 1px #ccc;

}

.add-client-form input {
    display: block;
    width: 95%;
    padding: 8px;
    margin: 10px auto;
    border: 1px solid #ccc;
    border-radius: 5px;
}

.add-client-form button {
    padding: 8px;
    margin: 10px auto;
    border: none;
    border-radius: 5px;
}

.add-client-form button {
    background-color: #5c0039;
    color: white;
    cursor: pointer;
}

.client-container {
    margin-left: -10px;
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
}

.client-card {
    background-color: #fff1f160;
    border-radius: 10px;
    padding: 20px;
    margin: 10px;
    width: 430px;
    transition: transform 0.3s ease;
    border: solid 1px #ccc;

}

.client-card:hover {
    transform: translateY(-5px);
}

.car-container {
    margin-top: 10px;

    border-radius: 5px;
}

.car-list {
    list-style: none;
    padding-left: 0;
    margin: 0;
    /* Removes default margin */
}

.car-item {
    display: flex;
    align-items: center;
    justify-content: flex-start;
    /* Aligns items to the start */
    margin-bottom: 5px;
    font-size: 0.9rem;
    /* Adjust font size if necessary */
    background-color: rgb(255, 255, 255);
    padding: 5px;
    border-radius: 10px;
    border: solid 1px #ccc;

}

.car-item span {
    display: flex;
    align-items: center;
    margin-right: 8px;
    /* Adjust space between elements */
    background-color: rgba(250, 235, 215, 0.274);
    padding: 2px 6px;
    font-size: small;
    border-radius: 5px;
}

i {
    margin-right: 10px;
    /* Adjust space between icon and text */
}

.car-item i {
    width: 20px;
    /* Ensures icons have the same width */
    text-align: center;
    /* Centers the icons */
}
</style>


<script>
import { API_URL } from '@/api';
export default {
    data() {
        return {
            newClient: {
                FirstName: '',
                LastName: '',
                PhoneNumber: ''
            },
            clients: [],
            cars: []
        };
    },
    created() {
        this.fetchClients();
        this.fetchCars();
    },
    methods: {
        addClient() {
            // Add your method for submitting the new client data to your server
            console.log('Adding client:', this.newClient);
            fetch(`${API_URL}/client`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(this.newClient)
            })
                .then(response => response.json())
                .then(data => {
                    console.log('Client added:', data);
                    this.fetchClients();
                })
                .catch(error => console.error('Error adding client:', error));

            this.fetchClients();
            this.newClient = { FirstName: '', LastName: '', PhoneNumber: '' };
        },
        fetchClients() {
            fetch(`${API_URL}/clients`)
                .then(response => response.json())
                .then(data => {
                    this.clients = data;
                })
                .catch(error => console.error('Error fetching clients:', error));
        },
        fetchCars() {
            fetch(`${API_URL}/carclient`)
                .then(response => response.json())
                .then(data => {
                    this.cars = data;
                })
                .catch(error => console.error('Error fetching cars:', error));
        },
        getCarsForClient(clientID) {
            return this.cars.filter(car => car.ClientID === clientID);
        },
        formatDate(date) {
            return new Date(date).toLocaleDateString('en-US', {
                year: 'numeric', month: 'long', day: 'numeric'
            });
        }
    }
}

</script>