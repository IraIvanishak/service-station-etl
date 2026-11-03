<template>
    <h1>Repairs</h1>
    <!-- Add button -->
    <div class="add-button" @click="openModal">+</div>

    <!-- Add Modal Component -->
    <AddStationModal v-if="isModalOpen" @close="isModalOpen = false" />
    <div class="container">
        <div class="repair-record" v-for="repair in repairs" :key="repair.repairId">
            <div :class="repair.repairStatus != 'Completed' ? 'header' : 'header-status'">
                <span class="repair-number">{{ repair.repairNumber }}</span>
                <span class="total-sum">${{ repair.totalSum.toFixed(2) }}</span>
                <div v-if="repair.totalSum === 0">
                    <span class="repair-status"> no services attached to order</span>
                </div>
                <div v-else>
                    <span class="repair-status">{{ repair.repairStatus }} - {{ repair.paymentStatus }}</span>
                </div>
                <div class="edit-delete-buttons">
                    <span v-if="repair.repairStatus != 'Completed'" class="edit-icon"
                        @click="editRepair(repair.repairId)">✎</span>
                    <span class="delete-icon" @click="deleteRepair(repair.repairId)">❌</span>
                </div>

            </div>
            <div class="content">
                <div class="notes">{{ repair.notes }}</div>
                <div class="vehicle-info">
                    <div>{{ repair.vehicleInfo.model }}</div>
                    <div>{{ repair.vehicleInfo.year }}</div>
                    <div>{{ repair.vehicleInfo.registrationNumber }}</div>
                </div>
                <div class="client-info">
                    <h5>client</h5>
                    <div>{{ repair.clientName }}</div>
                    <h5>worker</h5>
                    <div>{{ repair.assignedEmployee }}</div>
                </div>
                <div v-if="repair.totalSum > 0" class="client-info">
                    <h5>start date</h5>
                    <div>{{ new Date(repair.startDate).toLocaleDateString() }}</div>
                </div>


                <div v-if="repair.repairStatus != 'Completed'" class="add-service-button" @click="showServiceDrop(repair.repairId)">+</div>
                <br>
                <div v-if="showServiceDropdown && addingServiceToRepair === repair.repairId" class="service-dropdown">
                    <label for="name">Service</label>
                    <select id="name" v-model="newService.serviceId">
                        <option v-for="service in servicesList" :value="service.ServiceID">{{ service.Name }}
                        </option>
                    </select>
                    <label for="number">X</label>

                    <input style="width: 80px;" id="number" type="number" v-model="newService.quantity"
                        placeholder="Qty" min="1">
                    <button class="add" @click="addServiceToRepair(repair.repairId)">Add</button>
                </div>

                <div v-if="repair.totalSum > 0" class="services">
                    <div class="service" v-for="service in repair.services" :key="service.serviceName">

                        <div class="service-name">{{ service.serviceName }} {{ service.totalPrice / service.quantity }}
                            x
                            {{ service.quantity }}</div>
                        <div>{{ new Date(service.completionDate).toLocaleDateString() }}</div>
                        <div>Total for services: ${{ service.totalPrice }}</div>
                        <div class="service-controls">

                        </div>
                        <!-- Only render this div if there are inventory items -->
                        <div v-if="service.inventoryUsed && service.inventoryUsed.length "
                            v-for="item in service.inventoryUsed" :key="item.itemDescription">

                            <div v-if="item.amount != 0" class="inventory">
                                <div>{{ item.itemDescription }}</div>
                                <div>${{ item.pricePerUnit }} x {{ item.amount }}</div>
                            </div>
                        </div>
                        <div v-if="repair.repairStatus != 'Completed'" class="add-service-button" @click="showDetailDrop(service.serviceName, repair.repairId)">+
                        </div>
                        <br>
                        <div v-if="showDetailDro && newRepair === repair.repairId && service.serviceName === newServiceName"
                            class="service-dropdown">
                            <label for="name">Detail</label>
                            <select id="name" v-model="newDetail.detailId">
                                <option v-for="detail in details" :value="detail.DetailID">
                                    {{ detail.Description }}
                                </option>
                            </select>
                            <label for="numberD">X</label>

                            <input style="width: 80px;" id="numberD" type="number" v-model="newDetail.quantity"
                                placeholder="Qty" min="1">
                            <button class="add" @click="adddDetailToService(repair.repairId, service.serviceName)">Add</button>
                        </div>

                    </div>
                </div>

            </div>
        </div>
    </div>
</template>

<style scoped>
input[type="number"] {
    width: 50px;
    padding: 14px;
    border: 1px solid #ccc;
    border-radius: 5px;
    margin-top: 0;
}
.add {
    background-color: rgb(71, 71, 71);
    color: #ffffff;
    border: none;
    padding: 5px 10px;
    border-radius: 5px;
    cursor: pointer;
    width: 100%;
}

.service-controls {
    display: flex;
    align-items: center;
    gap: 10px;
}

.add-service-button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 25px;
    height: 25px;
    border: 1px solid #f0c0cb;
    border-radius: 50%;
    color: #83142c;
    cursor: pointer;
    font-size: 20px;
    line-height: 0;
}

.service-dropdown {
    display: flex;
    height: 28px;
    gap: 10px;
}

.add-button {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 50px;
    height: 50px;
    border-radius: 25px;
    background-color: #ffc0cb;
    /* Light pink background */
    color: #83142c;
    /* Dark pink plus symbol */
    font-size: 30px;
    position: fixed;
    bottom: 20px;
    right: 20px;
    cursor: pointer;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
    /* Subtle drop shadow */
    z-index: 500;
    /* Make sure it's above other elements but below the modal */
}

h1 {
    color: #333;
}


h5 {
    margin: 0;
    padding: 0;
}

.container {
    -moz-column-count: 2;
    -webkit-column-count: 2;
    column-count: 2;
    -moz-column-gap: 20px;
    -webkit-column-gap: 20px;
    column-gap: 20px;
}

.repair-record {
    display: inline-block;
    /* This makes the block element flow with the column layout */
    width: 100%;
    /* Full width of the column */
    background-color: #fff5f5;
    border: 1px solid #f0c0cb;
    border-radius: 10px;
    margin-bottom: 20px;
    /* Spacing between items */
    break-inside: avoid;
    /* Prevents items from splitting across columns */
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

@media (max-width: 768px) {
    .container {
        -moz-column-count: 1;
        -webkit-column-count: 1;
        column-count: 1;
    }
}

.header {
    background-color: #ffc0cb;
    color: #83142c;
    display: flex;
    justify-content: space-between;
    padding: 10px 20px;
    align-items: center;
}

.header .repair-number {
    font-weight: bold;
}

.header .total-sum {
    background-color: #fff;
    padding: 2px 5px;
    border-radius: 5px;
}

.header-status {
    background-color: #d3ffa9;
    color: black;
    display: flex;
    justify-content: space-between;
    padding: 10px 20px;
    align-items: center;
}

.header-status .repair-number {
    font-weight: bold;
}

.header-status .total-sum {
    background-color: #fff;
    padding: 2px 5px;
    border-radius: 5px;
}


.content {
    padding: 20px;
}

.notes {
    font-style: italic;
}

.vehicle-info {
    background-color: #fff0f1;
    padding: 10px;
    margin: 10px 0;
    border-radius: 5px;
}

.client-info {
    text-align: right;
}

.services {
    margin-top: 10px;
}

.service {
    background-color: #fff0f1;
    padding: 10px;
    margin-bottom: 10px;
    border-radius: 5px;
}

.service-name {
    font-weight: bold;
}

.inventory {
    background-color: #fff;
    margin: 5px 0;
    padding: 5px;
    border-radius: 5px;
    display: flex;
    justify-content: space-between;
}

.total-price {
    text-align: right;
    font-weight: bold;
}

.edit-delete-buttons {
    display: flex;
    justify-content: center;
    align-items: center;
    margin-top: 10px;
    scale: 1.5;
}

.edit-delete-buttons span {
    margin: 0 5px;
    cursor: pointer;
    font-size: 16px;
}

.edit-delete-buttons span:hover {
    color: #ff0000;
    /* Червоний колір при наведенні */
}

@media (max-width: 768px) {
    .container {
        padding: 10px;
    }

    .repair-record {
        width: 100%;
    }
}
</style>

<script>
import { API_URL } from '@/api';
import AddStationModal from './AddStationModal.vue';

export default {
    components: {
        AddStationModal
    },
    data() {
        return {
            showDetailDro: false,
            newServiceName: '',
            newRepair: '',
            newDetail: {
                detailId: null,
                quantity: 1
            },
            servicesList: [], // This will hold the list of services fetched from the API
            newService: {
                serviceId: null,
                quantity: 1
            },
            showServiceDropdown: false,
            repairs: [],
            isModalOpen: false,
            addingServiceToRepair: 0,
            details: []
        };
    },
    props: {
        id: {
            type: Number,
            required: true
        }
    },
    name: 'StationDetail',
    created() {
        this.fetchRepairs();
        this.fetchServices();
        this.getDetails();
    },
    methods: {
        async adddDetailToService(repairId, serviceName) {
            console.log(this.newDetail.detailId)
            try {
                console.log('Adding detail to service:', this.newDetail);
                const response = await fetch(`${API_URL}/repairsdetails`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        repairId: repairId,
                        serviceName: this.newServiceName,
                        detailId: this.newDetail.detailId,
                        quantity: this.newDetail.quantity
                    })
                });
            } catch (error) {
                console.error('Error adding detail to service:', error);
            }

            this.fetchRepairs();
            this.showDetailDro = false;
        },
        showDetailDrop(serviceName, repair) {
            this.showDetailDro = true;
            this.newServiceName = serviceName;
            this.newRepair = repair
        },
        async getDetails() {
            try {
                const response = await fetch(`${API_URL}/details`);
                if (response.ok) {
                    const details = await response.json();
                    this.details = details;
                } else {
                    console.error('Failed to fetch details.');
                }
            } catch (error) {
                console.error('Error fetching details:', error);
            }
        },
        async addServiceToRepair() {
            try {
                console.log('Adding service to repair:', this.newService);
                const response = await fetch(`${API_URL}/repairs/${this.addingServiceToRepair}/services`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(this.newService)
                });
                if (response.ok) {
                    console.log('Service added to repair successfully.');
                    this.fetchRepairs();
                } else {
                    console.error('Failed to add service to repair.');
                }
            } catch (error) {
                console.error('Error adding service to repair:', error);
            }
        },
        showServiceDrop(repairId) {
            this.showServiceDropdown = !this.showServiceDropdown;
            console.log('show service dropdown');
            this.addingServiceToRepair = repairId;
        },
        async fetchServices() {
            try {
                const response = await fetch(`${API_URL}/services`);
                if (response.ok) {
                    const services = await response.json();
                    this.servicesList = services;
                } else {
                    console.error('Failed to fetch services.');
                }
            } catch (error) {
                console.error('Error fetching services:', error);
            }
        },
        openModal() {
            this.isModalOpen = true;
            console.log('open modal');
        },
        async editRepair(repairId) {
            console.log('edit repair with ID:', repairId);
        },
        async fetchRepairs() {
            const stationId = this.$route.params.id;
            try {
                const response = await fetch(`${API_URL}/repairs`);
                const repairs = await response.json();
                this.repairs = repairs;
                this.repairs = this.repairs.sort((a, b) => b.repairId - a.repairId);
                console.log(this.repairs);
            } catch (error) {
                console.error(error);
            }

        },
        async deleteRepair(repairId) {
            try {
                await fetch(`${API_URL}/repairs/${repairId}`, {
                    method: 'DELETE'
                });
                // Optionally, remove the deleted repair from the local array
                this.repairs = this.repairs.filter(repair => repair.repairId !== repairId);
                console.log(`Repair with ID ${repairId} deleted successfully.`);
                this.fetchRepairs();
            } catch (error) {
                console.error(`Error deleting repair with ID ${repairId}:`, error);
            }
        },
    },
    mounted() {
        console.log('Component mounted.');
        console.log(this.id)
    },

}
</script>