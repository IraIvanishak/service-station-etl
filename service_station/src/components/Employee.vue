<template>
    <h2>Employee</h2>
    <div class="filter-panel">
        <div>
            <select v-model="filters.role">
                <option value="">Select Role</option>
                <option v-for="role in uniqueRoles" :key="role" :value="role">{{ role }}</option>
            </select>
            <select v-model="filters.station">
                <option value="">Select Station</option>
                <option v-for="station in uniqueStations" :key="station" :value="station">{{ station }}</option>
            </select>
        </div>
    </div>
    <div class="employee-container">
        <div v-for="employee in filteredEmployees" :key="employee.EmployeeID" class="employee-card">
            <h2><i class="fas fa-user-tie"></i> {{ employee.FirstName }} {{ employee.LastName }}</h2>
            <p><i class="fas fa-briefcase"></i> {{ employee.Role.String }}</p>
            <p><i class="fas fa-map-marker-alt"></i> {{ employee.ServiceStation.String }}</p>
        </div>
    </div>
</template>
<style scoped>

.filter-panel {
    background-color: #fcc4dfc5;
    padding: 10px;
    margin: 6px;
    margin-right: 30px;
    border-radius: 10px;
    margin-bottom: 10px;
}

.filter-panel div {
    display: flex;
    align-items: center;
    gap: 10px;
}

.filter-panel select, .filter-panel button {
    padding: 8px;
    border-radius: 5px;
    border: 1px solid #ccc;
    background-color: white;
}

.filter-panel button {
    cursor: pointer;
    background-color: #5c0039;
    color: white;
}
.employee-container {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-around;
    padding: 20px;
    margin-left: -30px;
}

.employee-card {
    border: 1px solid #ffd9e5;

    background-color: #fff;
    border-radius: 10px;
    padding: 20px;
    margin: 10px;
    width: 300px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.027);
    transition: transform 0.3s ease;
}

.employee-card:hover {
    transform: translateY(-5px);
}

i {
    margin-right: 10px;
}
</style>

<script>
import { API_URL } from '@/api';
export default {
    data() {
        return {
            employees: [],
            filters: {
                role: '',
                station: ''
            },
            uniqueRoles: [],
            uniqueStations: []
        };
    },
    created() {
        this.fetchEmployees();
    },
    computed: {
        filteredEmployees() {
            return this.employees.filter(employee => {
                return (!this.filters.role || employee.Role.String === this.filters.role) &&
                       (!this.filters.station || employee.ServiceStation.String === this.filters.station);
            });
        }
    },
    methods: {
        fetchEmployees() {
            fetch(`${API_URL}/stations`)
                .then(response => response.json())
                .then(data => {
                    this.uniqueStations = data.map(station => station.Address);
                })
                .catch(error => console.error('Error fetching stations:', error));

            fetch(`${API_URL}/roles`)
                .then(response => response.json())
                .then(data => {
                    this.uniqueRoles = data.map(role => role.Name);
                })
                .catch(error => console.error('Error fetching roles:', error));

            fetch(`${API_URL}/employees-full`)
                .then(response => response.json())
                .then(data => {
                    this.employees = data;
                })
                .catch(error => console.error('Error fetching employees:', error));
        },
        applyFilters() {
            // The computed property handles the filtering
        },
        clearFilters() {
            this.filters.role = '';
            this.filters.station = '';
        }
    }
}
</script>
