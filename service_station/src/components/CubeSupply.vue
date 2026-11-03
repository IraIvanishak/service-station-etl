<template>
    <h1>Service Inventory Cube</h1>
    <div class="service-inventory">
        <!-- Dimension selection -->
        <div class="grid-container">
            <!-- Left column for dimensions selection -->
            <div class="column">
                <fieldset class="lay">
                    <legend>Choose metrics to display:</legend>
                    <select v-model="selectedDimensions" multiple @change="updateActiveColumns">
                        <option v-for="(value, key) in dimensions" :key="key" :value="key">
                            {{ key }}
                        </option>
                    </select>
                </fieldset>

            </div>

            <!-- Right column for metrics checkboxes -->
            <div class="column">
                <fieldset class="lay">
                    <legend>Choose metrics to display:</legend>
                    <div>
                        <input type="checkbox" id="partsCostPercentage" v-model="TotalCost">
                        <label for="partsCostPercentage">TotalCost</label>
                    </div>
                    <div>
                        <input type="checkbox" id="serviceDurationDeviation" v-model="Quantity">
                        <label for="serviceDurationDeviation">Quantity</label>
                    </div>
                </fieldset>
            </div>
        </div>
        <!-- Filters container -->
        <div class="form-section">
            <!-- Filters for each selected dimension -->
            <div v-for="dimension in selectedDimensions" :key="dimension" class="form-row">
                <div>
                    <h3>{{ dimension }}</h3>
                    <!-- Dropdown to select specific filter value -->
                    <div  class="select-container">
                        <select v-model="selectedFilters[keyColumns[dimension]]" class="cselect">
                            <option disabled value="">{{ 'Select ' + keyColumns[dimension] }}</option>
                            <option :value="null">ALL</option> <!-- Reset option -->
                            <option v-for="value in uniqueFilters[keyColumns[dimension]]" :key="value" :value="value">
                                {{ value }}
                            </option>
                        </select>
                    </div>

                </div>
            </div>
        </div>


        <!-- Data table -->
        <table>
            <thead>
                <tr>
                    <th v-for="column in activeColumns" :key="column">{{ column }}</th>
                </tr>
            </thead>
            <tbody>
                <!-- Check if items are present -->
                <tr v-if="applyFilters().length > 0" v-for="item in applyFilters()" :key="item.ServiceInventoryFactID">
                    <td v-for="field in activeFieldNames" :key="item.ServiceInventoryFactID + field">{{ item[field] }}
                    </td>
                </tr>
                <!-- Message when no items are found -->
                <tr v-else>
                    <td :colspan="activeColumns.length">No records found for these dimmenseion values.</td>
                </tr>
            </tbody>
        </table>

    </div>
</template>



<script>
import { API_URL } from '@/api';
import axios from 'axios';

export default {
    name: 'CubeSupply',
    data() {
        return {
            TotalCost: false,
            Quantity: false,

            sullyData: [],
            selectedDimensions: ['Date', 'Supplier', 'Inventory'],
            dimensions: {
                Date: ['DateRange'],
                Supplier: ['SupplierName'],
                Inventory: ['Inventory Description', 'Inventory Article', 'Inventory Price']
            },
            dimensionColumns: {
                Date: ['DateRange'],
                Supplier: ['SupplierName'],
                Inventory: ['Inventory', 'Article', 'Price']
            },
            keyColumns: {
                Date: 'DateRange',
                Supplier: 'SupplierName',
                Inventory: 'InventoryDescription'
            },
            selectedKeyColumn: {
                Date: '',
                Supplier: '',
                Inventory: ''
            },
            uniqueFilters: {
                DateRange: [],
                SupplierName: [],
                InventoryDescription: []
            },
            selectedFilters: {
                DateRange: '',
                SupplierName: '',
                InventoryDescription: ''
            },
            uniqueValues: {
                DateRange: [],
                SupplierName: [],
                InventoryDescription: []
            },

        };
    },
    watch: {
        selectedDimensions: {
            handler() {
                this.updateUniqueValues();
            },
            immediate: true,  // This ensures it runs immediately after component mount
            deep: true
        }
    },

    mounted() {
        this.fetchInventoryData();
    },
    computed: {
        activeColumns() {
            let columns = [].concat(...this.selectedDimensions.map(d => this.dimensions[d]));
            if (this.TotalCost) columns.push('TotalCost');
            if (this.Quantity) columns.push('Quantity');
            return columns;
        },
        activeFieldNames() {
            let fields = [].concat(...this.selectedDimensions.map(d => this.dimensionColumns[d]));
            if (this.TotalCost) fields.push('TotalCost');
            if (this.Quantity) fields.push('Quantity');
            return fields;
        },
    },

    methods: {
        updateUniqueFilters() {
            const data = this.sullyData;
            this.uniqueFilters.DateRange = [...new Set(data.map(item => item.DateRange).filter(Boolean))];
            this.uniqueFilters.SupplierName = [...new Set(data.map(item => item.SupplierName).filter(Boolean))];
            this.uniqueFilters.InventoryDescription = [...new Set(data.map(item => item.Inventory).filter(Boolean))];
            console.log("Unique Filters:", JSON.stringify(this.uniqueFilters.DateRange, null, 2)); // Debug
        },
        fetchInventoryData() {
            axios.get(`${API_URL}/supply-cube`)
                .then(response => {
                    this.sullyData = response.data;
                    this.updateUniqueFilters();
                })
                .catch(error => {
                    console.error('There was an error fetching the inventory data:', error);
                });


        },
        applyFilters() {
            let result = this.sullyData;
            if (this.selectedFilters.DateRange) {
                result = result.filter(item => item.DateRange === this.selectedFilters.DateRange);
            }
            if (this.selectedFilters.SupplierName) {
                result = result.filter(item => item.SupplierName === this.selectedFilters.SupplierName);
            }
            if (this.selectedFilters.InventoryDescription) {
                result = result.filter(item => item.Inventory === this.selectedFilters.InventoryDescription);
            }
            return result;
        },
        updateUniqueValues() {
            if (this.sullyData.length > 0) {
                this.uniqueValues.RepairNumber = [...new Set(this.sullyData.map(item => item.RepairNumber).filter(Boolean))];
                this.uniqueValues.InventoryDescription = [...new Set(this.sullyData.map(item => item.Inventory).filter(Boolean))];
                console.log("Unique Repair Numbers:", this.uniqueValues.RepairNumber); // Debug
            }
        }
    }
}
</script>

<style>
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
    margin: 10px;
}

.lay {
    height: 134px;
    background-color: white;
    border: 1px solid #cc5671;
}

select[multiple] {
    width: 100%;
    height: 120px;
    padding: 8px;
    overflow-y: auto;
}

/* Style for individual options */
select[multiple] option {
    padding: 6px;
    cursor: pointer;
    border-bottom: 1px solid #f0f0f0;
}

/* Hover effects for options */
select[multiple] option:hover {
    background-color: #f9f9f9;
}

/* Focus styles */
select[multiple]:focus {
    border-color: #bb00cc;
    /* Highlight border when focused */
    outline: none;
    /* Remove default focus outline */
}

p,
label {
    font:
        1rem 'Fira Sans',
        sans-serif;
}

input {
    margin: 0.4rem;
}

.grid-container {
    display: flex;
    justify-content: space-between;
    margin-bottom: 20px;
}

.column {
    flex: 1;
    margin-right: 40px;
}

.column:last-child {
    margin-right: 0;
}

select {
    width: 100%;
    border: none;
    padding: 0;
    margin: 0;
}

/* Загальні стилі */
.service-inventory {
    background-color: #fff5f5;
    border: 1px solid #f0c0cb;
    border-radius: 10px;
    margin-bottom: 20px;
    padding: 20px;
    break-inside: avoid;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}


h1 {
    color: #333;
}



/* Стилі для таблиці */
table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 20px;
    background-color: white;
    border-radius: 10px;
}

th,
td {
    border: 1px solid #ddd;
    padding: 8px;
    text-align: left;
}

th {
    background-color: #a1133e;
    color: #ffffff;
}

tr:nth-child(even) {
    background-color: #f9f9f9;
}

tr:hover {
    background-color: #f1f1f1;
}

/* Base container styling */
.filters-container {
    display: flex;
    flex-direction: column;
}

/* Individual filter dimension styling */
.dimension-filter {
    border-radius: 5px;
    flex: 1;
}

.dimension-filter h3 {
    margin-top: 0;
    color: #333;
    font-size: 1.2em;
}

/* Styling for select containers */
.select-container {
    margin-top: 10px;
}

/* General select styling */
.cselect {
    width: 100%;
    padding: 8px;
    border: 1px solid #ccc;
    border-radius: 4px;
    cursor: pointer;
}

.cselect:focus {
    border-color: #0056b3;
    outline: none;
}
</style>