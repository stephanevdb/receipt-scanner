package main

import (
	"database/sql"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func InitCollections(app *pocketbase.PocketBase) error {
	dao := app.Dao()
	// --- Receipts collection ---
	receiptsCollection, err := dao.FindCollectionByNameOrId("receipts")
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Creating 'receipts' collection...")
			receiptsCollection = &models.Collection{
				Name: "receipts",
				Type: "base",
				Schema: schema.NewSchema(
					&schema.SchemaField{
						Name:     "filename",
						Type:     schema.FieldTypeText,
						Required: true,
						Options:  &schema.TextOptions{Max: types.Pointer(255)},
					},
					&schema.SchemaField{
						Name:     "title",
						Type:     schema.FieldTypeText,
						Required: false,
						Options:  &schema.TextOptions{Max: types.Pointer(255)},
					},
					&schema.SchemaField{
						Name:     "date",
						Type:     schema.FieldTypeDate,
						Required: false,
					},
					&schema.SchemaField{
						Name:     "total",
						Type:     schema.FieldTypeNumber,
						Required: false,
					},
					&schema.SchemaField{
						Name:     "user",
						Type:     schema.FieldTypeRelation,
						Required: true,
						Options:  &schema.RelationOptions{CollectionId: "_pb_users_auth_", MaxSelect: types.Pointer(1), CascadeDelete: false},
					},
				),
			}
			if err = dao.SaveCollection(receiptsCollection); err != nil {
				log.Printf("Error creating 'receipts' collection: %v", err)
				return err
			}
			log.Println("'receipts' collection created successfully.")
		} else {
			log.Printf("Error finding 'receipts' collection: %v", err)
			return err
		}
	} else {
		log.Println("'receipts' collection already exists.")
		// Check if the 'filename' field exists and add it if not
		if receiptsCollection.Schema.GetFieldByName("filename") == nil {
			log.Println("Adding 'filename' field to 'receipts' collection...")
			field := &schema.SchemaField{
				Name:     "filename",
				Type:     schema.FieldTypeText,
				Required: true,
				Options:  &schema.TextOptions{Max: types.Pointer(255)},
			}
			receiptsCollection.Schema.AddField(field)
			if err := dao.SaveCollection(receiptsCollection); err != nil {
				log.Printf("Error adding 'filename' field to 'receipts' collection: %v", err)
				return err
			}
			log.Println("'filename' field added successfully.")
		}
		if receiptsCollection.Schema.GetFieldByName("total") == nil {
			log.Println("Adding 'total' field to 'receipts' collection...")
			field := &schema.SchemaField{
				Name:     "total",
				Type:     schema.FieldTypeNumber,
				Required: false,
			}
			receiptsCollection.Schema.AddField(field)
			if err := dao.SaveCollection(receiptsCollection); err != nil {
				log.Printf("Error adding 'total' field to 'receipts' collection: %v", err)
				return err
			}
			log.Println("'total' field added successfully.")
		}
		if receiptsCollection.Schema.GetFieldByName("verified_total") == nil {
			log.Println("Adding 'verified_total' field to 'receipts' collection...")
			field := &schema.SchemaField{
				Name:     "verified_total",
				Type:     schema.FieldTypeBool,
				Required: false,
			}
			receiptsCollection.Schema.AddField(field)
			if err := dao.SaveCollection(receiptsCollection); err != nil {
				log.Printf("Error adding 'verified_total' field to 'receipts' collection: %v", err)
				return err
			}
			log.Println("'verified_total' field added successfully.")
		}
	}

	// --- Items collection ---
	itemsCollection, err := dao.FindCollectionByNameOrId("items")
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Creating 'items' collection...")
			itemsCollection := &models.Collection{
				Name: "items",
				Type: "base",
				Schema: schema.NewSchema(
					&schema.SchemaField{
						Name:     "name",
						Type:     schema.FieldTypeText,
						Required: true,
						Options:  &schema.TextOptions{Max: types.Pointer(255)},
					},
					&schema.SchemaField{
						Name:     "price",
						Type:     schema.FieldTypeNumber,
						Required: true,
					},
					&schema.SchemaField{
						Name:     "quantity",
						Type:     schema.FieldTypeNumber,
						Required: false,
					},
					&schema.SchemaField{
						Name:     "amount",
						Type:     schema.FieldTypeNumber,
						Required: false,
					},
					&schema.SchemaField{
						Name:     "paid",
						Type:     schema.FieldTypeNumber,
						Required: false,
					},
					&schema.SchemaField{
						Name:     "receipt",
						Type:     schema.FieldTypeRelation,
						Required: true,
						Options:  &schema.RelationOptions{CollectionId: receiptsCollection.Id, MaxSelect: types.Pointer(1), CascadeDelete: true},
					},
				),
			}
			if err = dao.SaveCollection(itemsCollection); err != nil {
				log.Printf("Error creating 'items' collection: %v", err)
				return err
			}
			log.Println("'items' collection created successfully.")
		} else {
			log.Printf("Error finding 'items' collection: %v", err)
			return err
		}
	} else {
		log.Println("'items' collection already exists.")
		// Ensure the relation is correctly pointing to the receipts collection
		receiptRelation := itemsCollection.Schema.GetFieldByName("receipt")
		if receiptRelation != nil {
			options, ok := receiptRelation.Options.(*schema.RelationOptions)
			if ok && options.CollectionId != receiptsCollection.Id {
				log.Println("Updating 'receipt' relation in 'items' collection...")
				options.CollectionId = receiptsCollection.Id
				if err := dao.SaveCollection(itemsCollection); err != nil {
					log.Printf("Error updating 'receipt' relation: %v", err)
					return err
				}
				log.Println("'receipt' relation updated successfully.")
			}
		}
		if itemsCollection.Schema.GetFieldByName("quantity") == nil {
			log.Println("Adding 'quantity' field to 'items' collection...")
			field := &schema.SchemaField{
				Name:     "quantity",
				Type:     schema.FieldTypeNumber,
				Required: false,
			}
			itemsCollection.Schema.AddField(field)
			if err := dao.SaveCollection(itemsCollection); err != nil {
				log.Printf("Error adding 'quantity' field to 'items' collection: %v", err)
				return err
			}
			log.Println("'quantity' field added successfully.")
		}
		if itemsCollection.Schema.GetFieldByName("amount") == nil {
			log.Println("Adding 'amount' field to 'items' collection...")
			field := &schema.SchemaField{
				Name:     "amount",
				Type:     schema.FieldTypeNumber,
				Required: false,
			}
			itemsCollection.Schema.AddField(field)
			if err := dao.SaveCollection(itemsCollection); err != nil {
				log.Printf("Error adding 'amount' field to 'items' collection: %v", err)
				return err
			}
			log.Println("'amount' field added successfully.")
		}
		if itemsCollection.Schema.GetFieldByName("paid") == nil {
			log.Println("Adding 'paid' field to 'items' collection...")
			field := &schema.SchemaField{
				Name:     "paid",
				Type:     schema.FieldTypeNumber,
				Required: false,
			}
			itemsCollection.Schema.AddField(field)
			if err := dao.SaveCollection(itemsCollection); err != nil {
				log.Printf("Error adding 'paid' field to 'items' collection: %v", err)
				return err
			}
			log.Println("'paid' field added successfully.")
		}
	}

	// --- Friend Groups collection ---
	friendGroupsCollection, err := dao.FindCollectionByNameOrId("friend_groups")
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Creating 'friend_groups' collection...")
			friendGroupsCollection = &models.Collection{
				Name:       "friend_groups",
				Type:       "base",
				ViewRule:   types.Pointer("@request.auth.id != ''"),
				CreateRule: types.Pointer("@request.auth.id != ''"),
				Schema: schema.NewSchema(
					&schema.SchemaField{
						Name:     "name",
						Type:     schema.FieldTypeText,
						Required: true,
						Options:  &schema.TextOptions{Max: types.Pointer(255)},
					},
				),
			}
			if err = dao.SaveCollection(friendGroupsCollection); err != nil {
				log.Printf("Error creating 'friend_groups' collection: %v", err)
				return err
			}
			log.Println("'friend_groups' collection created successfully.")
		} else {
			log.Printf("Error finding 'friend_groups' collection: %v", err)
			return err
		}
	} else {
		log.Println("'friend_groups' collection already exists.")
	}

	// --- Users collection ---
	// Add "name" field to users collection if it doesn't exist
	usersCollection, err := dao.FindCollectionByNameOrId("users")
	if err != nil {
		log.Printf("Could not find 'users' collection to add 'name' field, this might be expected if it's not created yet: %v", err)
	} else {
		// This ensures a user can update their own record, which is the default for
		// user collections but is good to set explicitly when modifying the schema.
		if usersCollection.UpdateRule == nil || *usersCollection.UpdateRule == "" {
			usersCollection.UpdateRule = types.Pointer("@request.auth.id = id")
			if err := dao.SaveCollection(usersCollection); err != nil {
				log.Printf("Error setting UpdateRule on 'users' collection: %v", err)
				return err
			}
			log.Println("UpdateRule for 'users' collection set successfully.")
		}

		// Allow using email as username to prevent validation errors on update
		if val, _ := usersCollection.Options["allowEmailUsername"].(bool); !val {
			log.Println("Enabling 'AllowEmailUsername' for 'users' collection...")
			usersCollection.Options["allowEmailUsername"] = true
			if err := dao.SaveCollection(usersCollection); err != nil {
				log.Printf("Error enabling 'AllowEmailUsername' on 'users' collection: %v", err)
				return err
			}
			log.Println("'AllowEmailUsername' for 'users' collection enabled successfully.")
		}

		if usersCollection.Schema.GetFieldByName("name") == nil {
			log.Println("Adding 'name' field to 'users' collection...")
			field := &schema.SchemaField{
				Name:     "name",
				Type:     schema.FieldTypeText,
				Required: false,
				Options:  &schema.TextOptions{Max: types.Pointer(255)},
			}
			usersCollection.Schema.AddField(field)
			if err := dao.SaveCollection(usersCollection); err != nil {
				log.Printf("Error adding 'name' field to 'users' collection: %v", err)
				return err
			}
			log.Println("'name' field added successfully.")
		}
		if usersCollection.Schema.GetFieldByName("iban") == nil {
			log.Println("Adding 'iban' field to 'users' collection...")
			field := &schema.SchemaField{
				Name:     "iban",
				Type:     schema.FieldTypeText,
				Required: false,
				Options:  &schema.TextOptions{Max: types.Pointer(34)},
			}
			usersCollection.Schema.AddField(field)
			if err := dao.SaveCollection(usersCollection); err != nil {
				log.Printf("Error adding 'iban' field to 'users' collection: %v", err)
				return err
			}
			log.Println("'iban' field added successfully.")
		}
		if usersCollection.Schema.GetFieldByName("bic") == nil {
			log.Println("Adding 'bic' field to 'users' collection...")
			field := &schema.SchemaField{
				Name:     "bic",
				Type:     schema.FieldTypeText,
				Required: false,
				Options:  &schema.TextOptions{Max: types.Pointer(11)},
			}
			usersCollection.Schema.AddField(field)
			if err := dao.SaveCollection(usersCollection); err != nil {
				log.Printf("Error adding 'bic' field to 'users' collection: %v", err)
				return err
			}
			log.Println("'bic' field added successfully.")
		}
		if usersCollection.Schema.GetFieldByName("friend_group") == nil {
			log.Println("Adding 'friend_group' relation to 'users' collection...")
			field := &schema.SchemaField{
				Name:     "friend_group",
				Type:     schema.FieldTypeRelation,
				Required: false,
				Options: &schema.RelationOptions{
					CollectionId:  friendGroupsCollection.Id,
					MaxSelect:     types.Pointer(1),
					CascadeDelete: false,
				},
			}
			usersCollection.Schema.AddField(field)
			if err := dao.SaveCollection(usersCollection); err != nil {
				log.Printf("Error adding 'friend_group' relation to 'users' collection: %v", err)
				return err
			}
			log.Println("'friend_group' relation added successfully.")
		}
	}

	return nil
}
