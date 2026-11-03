<template>
  <div class="modal-overlay" v-show="showModal" @click.self="closeModal">
     <div class="modal">
        <div class="modal-content">
           <h2>Add New Repair</h2>
           <div class="form-row">
              <label for="worker-select">Worker</label>
              <select id="worker-select" v-model="newRepair.assignedEmployeeId">
                 <option v-for="employee in employees" :key="employee.EmployeeID" :value="employee.EmployeeID">
                    {{ employee.FirstName }} {{ employee.LastName }}
                 </option>
              </select>
           </div>
           <div class="form-section">
              <div class="form-row">
                 <label for="client-select">Client</label>
                 <select id="client-select" v-model="newRepair.clientId">
                    <option v-for="client in clients" :key="client.ClientID" :value="client.ClientID">
                       {{ client.FirstName }} {{ client.LastName }}
                    </option>
                 </select>
              </div>
              <div class="form-row">
                 <label for="car-select">Car</label>
                 <select id="car-select" v-model="newRepair.carId">
                    <option v-for="car in cars" :key="car.CarID" :value="car.CarID">
                       {{ car.BrandName.String }} {{ car.Model.String }} - {{ car.Year.Int32 }}
                    </option>
                 </select>
              </div>
           </div>
           <div class="form-row">
              <label for="description">Description</label>
              <input id="description" type="text" v-model="newRepair.description" placeholder="Description">
           </div>
           <div class="form-section">
              <div class="form-row">
                 <label for="registration-number">Registration Number</label>
                 <input id="registration-number" type="text" v-model="newRepair.registrationNumber"
                    placeholder="Registration Number">
              </div>
              <div class="form-row">
                 <label for="vin">VIN</label>
                 <input id="vin" type="text" v-model="newRepair.vin" placeholder="Vehicle Identification Number">
              </div>
           </div>
            <div class="form-row">

           <label for="color">Color</label>
           <input id="color" type="text" v-model="newRepair.color" placeholder="Color">
           <!-- Input for notes -->
           <label for="notes">Notes</label>
           <input id="notes" type="text" v-model="newRepair.notes" placeholder="Notes">
           </div>
           <div class="form-section">
              <div class="form-row">
                 <!-- Dropdown for payment status -->
                 <select v-model="newRepair.paymentStatus">
                    <option value="1">Pending</option>
                    <option value="2">Paid</option>
                 </select>
              </div>
              <div class="form-row">
                 <select v-model="newRepair.repairStatus">
                    <option value="1">Scheduled</option>
                    <option value="2">In Progress</option>
                    <option value="3">Completed</option>
                 </select>
              </div>
           </div>
           <div class="form-section">
           <div class="form-row">
            <button class="cancel" @click="closeModal">Cancel</button>
            </div>
            <div class="form-row">
              <button @click="submitForm" type="submit">Submit</button>
           </div>
          </div>
        </div>
     </div>
  </div>
</template>


<script>
import { API_URL } from '@/api';
export default {
  name: 'AddStationModal',

  data() {
    return {
      showModal: true,
      employees: [], // Pre-fill with API call on mounted
      clients: [], // Pre-fill with API call on mounted
      cars: [], // Pre-fill with API call on mounted
      newRepair: {
        assignedEmployeeId: null,
        clientId: null,
        carId: null,
        notes: '',
        paymentStatus: '1', // Default to Pending
        repairStatus: '1', // Default to Scheduled
      }
    };
  },
  mounted() {
    this.fetchEmployees();
    this.fetchClients();
    this.fetchCars();
  },
  methods: {
    async fetchEmployees() {
      // API call to get employees and fill the this.employees array
      try {
        const response = await fetch(`${API_URL}/employees`);
        const data = await response.json();
        this.employees = data;
      } catch (error) {
        console.error('Error fetching employees:', error);
      }
    },
    async fetchClients() {
      // API call to get clients and fill the this.clients array
      try {
        const response = await fetch(`${API_URL}/clients`);
        const data = await response.json();
        this.clients = data;
      } catch (error) {
        console.error('Error fetching clients:', error);
      }
    },
    async fetchCars() {
      // API call to get cars and fill the this.cars array
      try {
        const response = await fetch(`${API_URL}/cars`);
        const data = await response.json();
        this.cars = data;
      } catch (error) {
        console.error('Error fetching cars:', error);
      }
    },
    closeModal() {
      this.showModal = false;
      this.$emit('close');
    },
    async submitForm() {
      // Implement form submission logic here
      console.log('Submitting new repair:', this.newRepair);

      console.log(JSON.stringify(this.newRepair));

      try {
        fetch(`${API_URL}/repairs`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(this.newRepair),
        });
      } catch (error) {
        console.error('Error submitting repair:', error);
      }


      // Perform the API call to create a new repair
      this.closeModal();
    }
  }
};
</script>

<style scoped>
.repair-form {
  display: flex;
  flex-direction: column;
}
.form-section {
    display: flex;
    justify-content: space-between;
    width: 100%;
}

.form-row {
  display: flex;
  flex-direction: column;
  flex-basis: 48%;
  /* Adjust based on gap */
  margin-bottom: 10px;
}


.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background-color: #fff;
  padding: 20px;
  border-radius: 10px;
  width: 80%;
  max-width: 500px;
  z-index: 1001;
}

input[type="text"],
select {
  height: 28px;
  width: 99%;
}

input,
select {
  width: 100%;
  /* This ensures the inputs fill the container */
  font-size: 1rem;
  line-height: 1.5;
  border-radius: 0.25rem;
  border: 1px solid #ced4da;
  /* Adjust color as needed */
}



button {
  z-index: 10;
  background-color: #640c26;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 10px 15px;
  cursor: pointer;
  margin-top: 10px;
  margin-right: 5px;
}
.cancel {
  z-index: 10;
  color: black;
  border: 1px solid #640c26;
  background-color: white;
  border-radius: 4px;
  padding: 10px 15px;
  cursor: pointer;
  margin-top: 10px;
  margin-right: 5px;
}
</style>