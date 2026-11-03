<style scoped>
.supply-container {
    background-color: #fff;
    border-radius: 5px;
    padding: 10px 20px;
    padding-left: 0;
}

.supply-header {
    color: #fff;
    background-color: #911247;
    /* Сірий фон для заголовку */
    border-radius: 5px;
    /* Закруглені кути */
    padding: 10px 20px;
    /* Внутрішній відступ */
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
    /* Відстань між заголовком та тілом */
}

.supply-body {
    background-color: #fff5f5;
    border: 1px solid #f0c0cb;
    border-radius: 10px;
    margin-bottom: 20px;
    padding: 20px;
    /* Spacing between items */
    break-inside: avoid;
    /* Prevents items from splitting across columns */
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.detail-item {
    border-bottom: 1px solid #e1e1e1;
    /* Лінія між деталями */
    padding: 10px 0;
    /* Відступи навколо деталі */
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.detail-description {
    /* Стиль для опису деталі */
    color: #333;
    /* Темний текст */
    flex: 1;
    /* Займає всю доступну ширину */
    margin-right: 20px;
    /* Відстань від ціни */
}

.detail-quantity,
.detail-price {
    /* Стилі для кількості та ціни */
    white-space: nowrap;
    /* Текст у один рядок */
}

.supplier-employee-section {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 10px;
    /* Відстань від деталей поставки */
}

.employee-name {
    font-style: italic;
    /* Курсив для імені співробітника */
    color: #555;
    /* Темно-сірий текст */
}

.status-tag {
    background-color: #ffffff;
    color: #2d572d;
    padding: 5px 10px;
    margin: 10px -5px;
    margin-top: 20px;
    margin-bottom: 10px;
    border-radius: 5px;
    /* Закруглені кути */
    font-size: 0.9em;
    /* Менший шрифт для статусу */
}

span {
    /* Стилі для всіх текстових елементів */
    margin-right: 10px;
    /* Відстань від сусідніх елементів */
}

.supplier-employee-section {
    flex-direction: column;
    align-items: flex-start;
}


.supply-form {
    display: flex;
    flex-direction: column;
    gap: 15px;
}

.supply-form h2 {
    color: #911247;
    margin-bottom: 10px;
}

.supply-form label {
    color: #333;
    font-size: 14px;
    margin-bottom: 5px;
}

.supply-form input[type="text"],
.supply-form input[type="date"] {
    border: 1px solid #ccc;
    border-radius: 4px;
    padding: 10px;
    font-size: 16px;
    color: #333;
}

.supply-form input[type="text"]:focus,
.supply-form input[type="date"]:focus {
    outline: none;
    border-color: #911247;
    box-shadow: 0 0 0 2px rgba(145, 18, 71, 0.2);
}

.supply-form button {
    width: 200px;
    padding: 10px 15px;
    background-color: #911247;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.3s ease;
    font-size: 16px;
}

.supply-form button:hover {
    background-color: #a12f5c;
}

.supply-form button:active {
    box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.2);
}

.fas {
    margin-right: 5px;
    /* додає відстань між іконкою і текстом */
    color: #911247;
    /* встановлює колір іконок, щоб відповідати дизайну */
}

.form-container {
    background-color: #fff5f5;
    border: 1px solid #f0c0cb;
    border-radius: 10px;
    margin: 20px;
    padding: 20px;
}

.supply-form {
    display: flex;
    justify-content: space-between;
    width: 100%;
}



.form-group {
    margin-bottom: 20px;
}

.detail-list {
    width: 100%;
    margin-top: 20px;
}

.detail-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
}

button {
    margin-top: 10px;
}

input,
select {
    font-size: 1rem;
    line-height: 1.5;
    border-radius: 0.25rem;
    border: 1px solid #ced4da;
    /* Adjust color as needed */
}

.supply-form-container {
    border: 1px solid #ced4da;

    margin: 20px;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 20px;
    background-color: #ffffff;
    border-radius: 10px;
}

.form-section {
    display: flex;
    justify-content: space-between;
    width: 100%;
    margin-bottom: 20px;
}

.form-row {
    display: flex;
    flex-direction: column;
    flex-basis: 48%;
    /* Adjust based on gap */
    margin-bottom: 10px;
}

label {
    margin-bottom: 5px;
    font-weight: bold;
    color: #333;
}

input[type="date"],
select {
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 4px;
    background: white;
}

.submit-button {
    padding: 10px 15px;
    background-color: #911247;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 16px;
}

.submit-button:hover {
    background-color: #a12f5c;
}

.detail-item-add {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
}

button {
    padding: 10px 15px;
    width: 100px;
    background-color: #911247;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 16px;

}


.set-delivered-btn {
    margin-top: 0px;
    margin-left: 20px;
    width: 120px;
    /* Button width */
    background-color: #D94662;
    /* Deep pink/red color */
    color: white;
    /* White text */
    border: none;
    /* No border */
    padding: 5px 20px;
    /* Top/bottom and left/right padding */
    font-size: 12px;
    /* Text size */
    border-radius: 5px;
    /* Rounded corners */
    cursor: pointer;
    /* Pointer cursor on hover */
    transition: background-color 0.3s;
    /* Smooth transition for hover effect */
}

.set-delivered-btn:hover {
    background-color: #C03250;
    /* Slightly darker shade for hover */
}

h1 {
    color: #333;
}

</style>

<template>
    <h1>Supply</h1>
    <AddSupply />
    <div>

        <main>
            <div class="supply-container" v-for="supply in supplies" :key="supply.supply_id">

                <div class="supply-header">
                    <span>Supply: {{ supply.supply_number }}</span>
                    <div class="edit-delete-buttons">
                        <span v-if="supply.supply_status === 'Pending'" class="edit-icon" @click="editSupply(supply)">✎</span>
                        <span class="delete-icon" @click="deleteSupply(supply.supply_id)">❌</span>
                    </div>
                </div>
                <div class="supply-body">
                    <div class="detail-item" v-for="detail in supply.details" :key="detail.detail_id">
                        <span class="detail-description">{{ detail.description }}</span>
                        <span class="detail-quantity"> {{ detail.quantity }} x </span>
                        <span class="detail-price">${{ detail.price_per_unit | currency }}</span>
                        <span class="invoice-number">{{ detail.invoice_number }}</span>
                        <span class="service-station">{{ detail.service_station }}</span>
                        <span>{{ new Date(detail.supply_date).toDateString() }}</span>

                    </div>
                    <div class="supplier-employee-section">
                        <span>
                            <i class="fas fa-truck"></i> 
                            Supplier: {{ supply.supplier_name }}
                        </span>
                        <span>
                            <i class="fas fa-user"></i> 
                            Handled by: {{ supply.employee.first_name }} {{ supply.employee.last_name }}
                        </span>
                    </div>
                    <div class="status-tag">
                        <i class="fas fa-check-circle"></i>
                        {{ supply.supply_status }}
                        <!-- Button to set status to delivered -->
                        <button class="set-delivered-btn" v-if="supply.supply_status === 'Pending'"
                            @click="setDelivered(supply.supply_id)">
                            Set Delivered
                        </button>
                    </div>
                </div>

            </div>
        </main>
    </div>
</template>

<script>
import { API_URL } from '@/api';
import AddSupply from "./AddSupply.vue"
export default {
    components: {
        AddSupply,
    },
    name: 'SupplyComponent',
    data() {
        return {
            supplies: [],
        };
    },
    mounted() {
        this.fetchSupplies();
    },
    methods: {
        async setDelivered(supplyId) {
            await fetch(`${API_URL}/supplies/${supplyId}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ status: 'Delivered' }),
            });
            try {
                const response = await fetch(`${API_URL}/supplies`);
                const supplies = await response.json();
                this.supplies = supplies.sort((a, b) => b.supply_id - a.supply_id);
            } catch (error) {
                console.error(error);
            }
        },
        removeDetail(index) {
            this.supplyForm.details.splice(index, 1);
        },
        async deleteSupply(supplyId) {
            try {
                await fetch(`${API_URL}/supplies/${supplyId}`, {
                    method: 'DELETE',
                });
                this.supplies = this.supplies.filter((supply) => supply.supply_id !== supplyId);
            } catch (error) {
                console.error(error);
            }
        },
        formatDate(value) {
            if (!value) return '';
            const date = new Date(value);
            return date.toLocaleDateString();
        },
        currency(value) {
            if (!value) return '';
            return `$${parseFloat(value).toFixed(2)}`;
        },
        async fetchSupplies() {
            try {
                const response = await fetch(`${API_URL}/supplies`);
                const supplies = await response.json();
                this.supplies = supplies.sort((a, b) => b.supply_id - a.supply_id);
            } catch (error) {
                console.error(error);
            }
        },
    }

}
</script>