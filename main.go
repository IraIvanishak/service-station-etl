package main

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/IraIvanishak/service-station-etl/dto"
	"github.com/IraIvanishak/service-station-etl/internal/config"
	"github.com/IraIvanishak/service-station-etl/storage/dbs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
)

var queries *dbs.Queries

func main() {
	cfg := config.Load()
	app := fiber.New()

	connection, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer connection.Close()
	if err := connection.Ping(); err != nil {
		log.Fatalf("Database is not reachable (DATABASE_URL): %v", err)
	}

	queries = dbs.New(connection)
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))

	app.Get("/city", func(ctx *fiber.Ctx) error {
		cities, err := queries.ListCities(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(cities)
	})

	app.Get("/stations", func(ctx *fiber.Ctx) error {
		stations, err := queries.ListServiceStations(context.Background())

		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(stations)
	})

	app.Post("/stations", func(ctx *fiber.Ctx) error {
		station := new(dto.Station)
		if err := ctx.BodyParser(station); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		ss := &dbs.AddServiceStationParams{
			Address: station.Address,
			CityID: sql.NullInt32{
				Int32: station.City,
				Valid: true,
			},
		}
		err := queries.AddServiceStation(context.Background(), *ss)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(station)
	})

	app.Delete("/stations/:id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		intid, err := strconv.Atoi(id)
		if err != nil {
			// If the ID is not an integer, return a Bad Request response.
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Station ID must be an integer",
			})
		}

		err = queries.DeleteServiceStation(context.Background(), int32(intid))
		if err != nil {
			// Handle the case where the station does not exist.
			if err == sql.ErrNoRows {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Station not found",
				})
			}
			// Handle other potential errors.
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.SendString("Deleted")
	})

	app.Put("/stations", func(ctx *fiber.Ctx) error {
		station := new(dto.Station)
		if err := ctx.BodyParser(station); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		ss := dbs.UpdateServiceStationParams{
			StationID: station.Id,
			Address:   station.Address,
			CityID: sql.NullInt32{
				Int32: station.City,
				Valid: true,
			},
		}
		err := queries.UpdateServiceStation(context.Background(), ss)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(station)
	})

	app.Get("/repairs", func(ctx *fiber.Ctx) error {
		repairs, err := queries.ListOrders(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		rep := dto.ConvertToDTO(repairs)
		for i, repair := range rep {
			rep[i].Services = dto.AggregateServices(repair.Services)
		}

		for i := 0; i < len(rep); i++ {
			sum := 0.0
			for _, service := range rep[i].Services {
				sum += service.TotalPrice
				for _, inventory := range service.InventoryUsed {
					sum += inventory.Amount * inventory.PricePerUnit
				}
			}
			rep[i].TotalSum = sum
		}
		return ctx.JSON(rep)
	})

	app.Get("/employees", func(ctx *fiber.Ctx) error {
		employees, err := queries.ListEmployees(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(employees)
	})

	app.Get("/clients", func(ctx *fiber.Ctx) error {
		clients, err := queries.ListClients(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(clients)
	})

	app.Get("/cars", func(ctx *fiber.Ctx) error {
		cars, err := queries.ListCars(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(cars)
	})

	app.Post("/repairs", func(ctx *fiber.Ctx) error {
		repair := new(dto.Repair)
		if err := ctx.BodyParser(repair); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// insert carclient
		ccid, err := queries.StoreCarClient(context.Background(), dbs.StoreCarClientParams{
			CarID:              repair.CarID,
			ClientID:           repair.ClientID,
			Description:        sql.NullString{String: repair.Description, Valid: true},
			RegistrationNumber: sql.NullString{String: repair.RegistrationNumber, Valid: true},
			Vin:                sql.NullString{String: repair.VIN, Valid: true},
			Color:              sql.NullString{String: repair.Color, Valid: true},
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = queries.StoreOrder(context.Background(), dbs.StoreOrderParams{
			EmployeeID:      sql.NullInt32{Int32: repair.AssignedEmployeeID, Valid: true},
			StartDate:       sql.NullTime{Time: time.Now(), Valid: true},
			RepairStatusID:  sql.NullInt32{Int32: 1, Valid: true},
			Notes:           sql.NullString{String: repair.Notes, Valid: true},
			RepairNumber:    sql.NullString{String: dto.CreateRepairNumber(), Valid: true},
			PaymentStatusID: sql.NullInt32{Int32: 1, Valid: true},
			CarClientID:     sql.NullInt32{Int32: ccid, Valid: true},
			TotalSum:        sql.NullString{String: "0", Valid: true},
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(repair)
	})

	app.Delete("/repairs/:id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		intid, err := strconv.Atoi(id)
		if err != nil {
			// If the ID is not an integer, return a Bad Request response.
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Repair ID must be an integer",
			})
		}

		err = queries.DeleteOrderCascade(context.Background(), int32(intid))
		if err != nil {
			// Handle the case where the repair does not exist.
			if err == sql.ErrNoRows {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Repair not found",
				})
			}
			// Handle other potential errors.
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.SendString("Deleted")

	})

	app.Get("/services", func(ctx *fiber.Ctx) error {
		services, err := queries.ListServices(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(services)
	})

	app.Post("/repairs/:id/services", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		intid, err := strconv.Atoi(id)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Repair ID must be an integer",
			})
		}

		service := new(dto.ServiceRepair)
		if err := ctx.BodyParser(service); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		price, _ := queries.GetServicePriceById(context.Background(), service.ServiceID)
		priceInt := dto.ParseFloat(price.String)
		err = queries.StoreServiceRepair(context.Background(), dbs.StoreServiceRepairParams{
			ServiceID: sql.NullInt32{Int32: service.ServiceID, Valid: true},
			Quantity:  sql.NullInt32{Int32: service.Quantity, Valid: true},
			RepairID:  sql.NullInt32{Int32: int32(intid), Valid: true},
			CompletionDate: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			ServiceTotalPrice: sql.NullString{
				String: strconv.FormatFloat(priceInt*float64(service.Quantity), 'f', 6, 64),
				Valid:  true,
			},
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(service)

	})

	app.Get("/supplies", func(ctx *fiber.Ctx) error {
		rows, err := queries.ListSupply(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		supplies, err := dto.GroupSupplyData(rows)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(supplies)
	})

	app.Get("/details", func(ctx *fiber.Ctx) error {
		rows, err := queries.ListDetails(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(rows)
	})

	app.Get("/suppliers", func(ctx *fiber.Ctx) error {
		rows, err := queries.ListSuppliers(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(rows)
	})

	app.Get("/meta", func(ctx *fiber.Ctx) error {
		meta, err := queries.MetaData(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(meta)
	})

	app.Post("/supplies", func(ctx *fiber.Ctx) error {
		dtoNS := new(dto.NewSupplyDTO)
		if err := ctx.BodyParser(dtoNS); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		stat := 2
		if dtoNS.Status == "Pending" {
			stat = 1
		}
		countTotalSum := func(details []dto.DetailsDTO) float64 {
			var sum float64
			for _, detail := range details {
				price, _ := queries.GetSupplyDetailPrice(context.Background(), detail.DetailID)
				priceInt := dto.ParseFloat(price.String)
				sum += float64(detail.Quantity) * priceInt
			}
			return sum
		}
		rows, _ := queries.ListSupply(context.Background())
		sunppliesNumber := len(rows)
		supplyID, err := queries.NewSupply(context.Background(), dbs.NewSupplyParams{
			TotalSum: sql.NullString{
				String: strconv.FormatFloat(countTotalSum(dtoNS.Details), 'f', 6, 64),
				Valid:  true,
			},
			StatusID: sql.NullInt32{
				Int32: int32(stat),
				Valid: true,
			},
			SupplyDate: dtoNS.SupplyDate,
			SupplierID: sql.NullInt32{
				Int32: dtoNS.SupplierID,
				Valid: true,
			},
			EmployeeID: sql.NullInt32{
				Int32: dtoNS.EmployeeID,
				Valid: true,
			},
			SupplyNumber: dto.CreateSupplyNumber(sunppliesNumber + 1),
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error queries.NewSupply update": err.Error(),
			})
		}

		if stat == 2 {
			_, err := queries.SetSupplyStatus(context.Background(), dbs.SetSupplyStatusParams{
				StatusID: sql.NullInt32{
					Int32: 2,
					Valid: true,
				},
				SupplyID: int32(supplyID),
			})
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		}
		for _, detail := range dtoNS.Details {
			// get detail price
			price, _ := queries.GetSupplyDetailPrice(context.Background(), detail.DetailID)
			priceInt := dto.ParseFloat(price.String)

			id, err := queries.NewSupplyDetail(context.Background(), dbs.NewSupplyDetailParams{
				Quantity: sql.NullInt32{
					Int32: detail.Quantity,
					Valid: true,
				},
				SupplyID: sql.NullInt32{
					Int32: supplyID,
					Valid: true,
				},
				DetailID: sql.NullInt32{
					Int32: detail.DetailID,
					Valid: true,
				},
				PricePerUnit: sql.NullString{
					String: strconv.FormatFloat(priceInt, 'f', 6, 64),
					Valid:  true,
				},
			})
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			idd, err := queries.StoreInventory(context.Background(), dbs.StoreInventoryParams{
				Quantity: sql.NullInt32{
					Int32: detail.Quantity,
					Valid: true,
				},
				DetailID: sql.NullInt32{
					Int32: detail.DetailID,
					Valid: true,
				},
			})
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			if stat == 2 {
				_, err = queries.NewSupplyInvoice(context.Background(), dbs.NewSupplyInvoiceParams{
					InventoryID: idd,
					SupplyDetailsID: sql.NullInt32{
						Int32: id,
						Valid: true,
					},
					InvoiceNumber: sql.NullString{
						String: dto.CreateInvoiceNumber(),
						Valid:  true,
					},
				})
				if err != nil {
					return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": err.Error(),
					})
				}
			}
		}

		return ctx.JSON(dtoNS)
	})

	app.Patch("/supplies/:id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		intid, err := strconv.Atoi(id)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Supply ID must be an integer",
			})
		}
		// supply, err := queries.GetSupplyByID(context.Background(), int32(intid))
		// if err != nil {
		// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 		"error": err.Error(),
		// 	})
		// }

		supplyDetails, err := queries.GetSupplyDetailsByID(context.Background(), sql.NullInt32{
			Int32: int32(intid),
			Valid: true,
		})

		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error GetSupplyDetailsByID	": err.Error(),
			})
		}

		for _, detail := range supplyDetails {
			idd, err := queries.GetInventoryIDbyDetailsID(context.Background(), detail.DetailID)

			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error GetInventoryIDbyDetailsID": err.Error(),
				})
			}

			_, err = queries.NewSupplyInvoice(context.Background(), dbs.NewSupplyInvoiceParams{
				InventoryID: idd,
				SupplyDetailsID: sql.NullInt32{
					Int32: detail.SupplyDetailsID,
					Valid: true,
				},
				InvoiceNumber: sql.NullString{
					String: dto.CreateInvoiceNumber(),
					Valid:  true,
				},
			})
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		}

		sipId, err := queries.SetSupplyStatus(context.Background(), dbs.SetSupplyStatusParams{
			StatusID: sql.NullInt32{
				Int32: 2,
				Valid: true,
			},
			SupplyID: int32(intid),
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.JSON("Updated supply with id: " + strconv.Itoa(int(intid)) + " and supply invoice id: " + strconv.Itoa(int(sipId)))
	})

	app.Delete("/supplies/:id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		intid, err := strconv.Atoi(id)
		if err != nil {
			// If the ID is not an integer, return a Bad Request response.
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Supply ID must be an integer",
			})
		}
		err = queries.DeleteSupplyDetailsBySupplyId(context.Background(), sql.NullInt32{
			Int32: int32(intid),
			Valid: true,
		})
		if err != nil {
			// Handle the case where the supply details do not exist.
			if err == sql.ErrNoRows {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Supply details not found",
				})
			}
			// Handle other potential errors.
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		err = queries.DeleteSupply(context.Background(), int32(intid))
		if err != nil {
			// Handle the case where the supply does not exist.
			if err == sql.ErrNoRows {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Supply not found",
				})
			}
			// Handle other potential errors.
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.SendString("Deleted")
	})

	app.Get("/service-invenroty-cube", func(ctx *fiber.Ctx) error {
		rows, err := queries.GetCubeServiceInventory(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		var userFriendly []dto.ServiceInventoryDTO
		for _, row := range rows {
			userFriendly = append(userFriendly, dto.ConvertToUserFriendly(row))
		}
		return ctx.JSON(userFriendly)
	})

	app.Get("/supply-cube", func(ctx *fiber.Ctx) error {
		rows, err := queries.GetCubeSupply(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var userFriendly []dto.SupplyDetailsDTO
		for _, row := range rows {
			userFriendly = append(userFriendly, dto.ToDTOSupply(row))
		}
		return ctx.JSON(userFriendly)
	})

	app.Post("/clearmeta", func(ctx *fiber.Ctx) error {
		// err := queries.ClearMeta(context.Background())
		// if err != nil {
		// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 		"error": err.Error(),
		// 	})
		// }

		err = queries.ClearStorage(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.SendString("Cleared")
	})

	app.Post("/init", func(ctx *fiber.Ctx) error {
		err := queries.OltpToStage(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error OltpToStage": err.Error(),
			})
		}

		err = queries.StageToOlap(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error StageToOlap": err.Error(),
			})
		}

		return ctx.SendString("Initialized")
	})

	app.Get("/stage", func(ctx *fiber.Ctx) error {
		err := queries.OlapToStageIncremental(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		stage, err := queries.GetStageData(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(stage)
	})

	app.Post("/incremental", func(ctx *fiber.Ctx) error {
		err := queries.OlapToStageIncremental(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		data, err := queries.GetStageData(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// count changed rows
		changed := 0
		for _, row := range data {
			changed += int(row.RowCount)
		}

		// count chabged tables
		tables := 0
		for _, row := range data {
			if row.RowCount > 0 {
				tables++
			}
		}

		err = queries.StageToOlapIncremental(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		dlid, err := queries.GetLastDataLoadHistoryID(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = queries.UpdateDataLoadHistory(context.Background(), dbs.UpdateDataLoadHistoryParams{
			LoadRows: sql.NullInt32{
				Int32: int32(changed),
				Valid: true,
			},
			AffectedTableCount: sql.NullInt32{
				Int32: int32(tables),
				Valid: true,
			},
			DataLoadHistoryID: int32(dlid),
		})

		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.SendString("Incremental")
	})

	app.Get("/carclient", func(ctx *fiber.Ctx) error {
		rows, err := queries.ListCarClients(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(rows)
	})

	app.Post("/client", func(ctx *fiber.Ctx) error {
		type Client struct {
			FirstName string `json:"FirstName"`
			LastName  string `json:"LastName"`
			Phone     string `json:"PhoneNumber"`
		}
		client := new(Client)
		if err := ctx.BodyParser(client); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		c := dbs.SoreClientParams{
			FirstName:   client.FirstName,
			LastName:    client.LastName,
			PhoneNumber: sql.NullString{String: client.Phone, Valid: true},
		}

		err := queries.SoreClient(context.Background(), c)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(client)
	})

	app.Get("/employees-full", func(ctx *fiber.Ctx) error {
		employees, err := queries.ListEmployees(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(employees)
	})

	app.Get("/roles", func(ctx *fiber.Ctx) error {
		roles, err := queries.ListEmployeesRoles(context.Background())
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(roles)
	})

	//{"repairId":1505,"serviceName":"Engine Belt Replacement","detailId":14,"quantity":2}
	app.Post("/repairsdetails", func(ctx *fiber.Ctx) error {
		type RepairService struct {
			RepairID    int32  `json:"repairId"`
			DetailId    int32  `json:"detailId"`
			ServiceName string `json:"serviceName"`
			Quantity    int32  `json:"quantity"`
		}
		repairService := new(RepairService)
		if err := ctx.BodyParser(repairService); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		sid, err := queries.GetServiceRepairByRepairIdAndServiceName(context.Background(), dbs.GetServiceRepairByRepairIdAndServiceNameParams{
			RepairID: sql.NullInt32{
				Int32: repairService.RepairID,
				Valid: true,
			},
			ServiceName: repairService.ServiceName,
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		iid, err := queries.GetInventoryIDbyDetailsID(context.Background(), sql.NullInt32{
			Int32: repairService.DetailId,
			Valid: true,
		})

		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error GetInventoryIDbyDetailsID": err.Error(),
			})
		}

		err = queries.StoreServiceInventory(context.Background(), dbs.StoreServiceInventoryParams{
			ServiceRepairID: sql.NullInt32{
				Int32: sid,
				Valid: true,
			},
			InventoryID: sql.NullInt32{
				Int32: iid,
				Valid: true,
			},
			Amount: sql.NullString{
				String: strconv.Itoa(int(repairService.Quantity)),
				Valid:  true,
			},
		})
		if err != nil {

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.JSON(repairService)
	})

	log.Fatal(app.Listen(cfg.Addr))
}
