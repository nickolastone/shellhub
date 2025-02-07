package mongo

import (
	"context"

	"github.com/shellhub-io/shellhub/api/pkg/gateway"
	"github.com/shellhub-io/shellhub/api/store/mongo/queries"
	"github.com/shellhub-io/shellhub/pkg/api/paginator"
	"github.com/shellhub-io/shellhub/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Store) UserList(ctx context.Context, pagination paginator.Query, filters []models.Filter) ([]models.User, int, error) {
	query := []bson.M{}

	if tenant := gateway.TenantFromContext(ctx); tenant != nil {
		query = append(query, bson.M{
			"$match": bson.M{
				"tenant_id": tenant.ID,
			},
		})
	}

	query = append(query, []bson.M{
		{
			"$addFields": bson.M{
				"user_id": bson.M{"$toString": "$_id"},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "namespaces",
				"localField":   "user_id",
				"foreignField": "owner",
				"as":           "namespaces",
			},
		},
		{
			"$addFields": bson.M{
				"namespaces": bson.M{"$size": "$namespaces"},
			},
		},
	}...)

	queryMatch, err := queries.BuildFilterQuery(filters)
	if err != nil {
		return nil, 0, fromMongoError(err)
	}

	if len(queryMatch) > 0 {
		query = append(query, queryMatch...)
	}

	queryCount := query
	queryCount = append(queryCount, bson.M{"$count": "count"})
	count, err := aggregateCount(ctx, s.db.Collection("users"), queryCount)
	if err != nil {
		return nil, 0, fromMongoError(err)
	}

	if pagination.Page > 0 && pagination.PerPage > 0 {
		query = append(query, queries.BuildPaginationQuery(pagination)...)
	}

	users := make([]models.User, 0)
	cursor, err := s.db.Collection("users").Aggregate(ctx, query)
	if err != nil {
		return nil, 0, fromMongoError(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		user := new(models.User)
		err = cursor.Decode(&user)
		if err != nil {
			return nil, 0, fromMongoError(err)
		}

		users = append(users, *user)
	}

	return users, count, fromMongoError(err)
}

func (s *Store) UserCreate(ctx context.Context, user *models.User) error {
	_, err := s.db.Collection("users").InsertOne(ctx, user)

	return fromMongoError(err)
}

func (s *Store) UserGetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := new(models.User)

	if err := s.db.Collection("users").FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		return nil, fromMongoError(err)
	}

	return user, nil
}

func (s *Store) UserGetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)

	if err := s.db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, fromMongoError(err)
	}

	return user, nil
}

func (s *Store) UserGetByID(ctx context.Context, id string, ns bool) (*models.User, int, error) {
	user := new(models.User)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, 0, err
	}

	if err := s.db.Collection("users").FindOne(ctx, bson.M{"_id": objID}).Decode(&user); err != nil {
		return nil, 0, fromMongoError(err)
	}

	if !ns {
		return user, 0, nil
	}

	nss := struct {
		NamespacesOwned int `bson:"namespacesOwned"`
	}{}

	query := []bson.M{
		{
			"$match": bson.M{
				"_id": objID,
			},
		},
		{
			"$addFields": bson.M{
				"_id": bson.M{
					"$toString": "$_id",
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "namespaces",
				"localField":   "_id",
				"foreignField": "owner",
				"as":           "ns",
			},
		},
		{
			"$addFields": bson.M{
				"namespacesOwned": bson.M{
					"$size": "$ns",
				},
			},
		},
		{
			"$project": bson.M{
				"namespacesOwned": 1,
				"_id":             0,
			},
		},
	}

	cursor, err := s.db.Collection("users").Aggregate(ctx, query)
	if err != nil {
		return nil, 0, fromMongoError(err)
	}

	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return nil, 0, fromMongoError(err)
	}

	if err = cursor.Decode(&nss); err != nil {
		return nil, 0, fromMongoError(err)
	}

	return user, nss.NamespacesOwned, nil
}

func (s *Store) UserUpdateData(ctx context.Context, id string, data models.User) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fromMongoError(err)
	}

	if _, err := s.db.Collection("users").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"name": data.Name, "username": data.Username, "email": data.Email, "last_login": data.LastLogin}}); err != nil {
		return fromMongoError(err)
	}

	return nil
}

func (s *Store) UserUpdatePassword(ctx context.Context, newPassword string, id string) error {
	if _, _, err := s.UserGetByID(ctx, id, false); err != nil {
		return fromMongoError(err)
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fromMongoError(err)
	}

	if _, err := s.db.Collection("users").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"password": newPassword}}); err != nil {
		return fromMongoError(err)
	}

	return nil
}

func (s *Store) UserUpdateFromAdmin(ctx context.Context, name string, username string, email string, password string, id string) error {
	user, _, err := s.UserGetByID(ctx, id, false)
	objID, _ := primitive.ObjectIDFromHex(id)

	if err != nil {
		return fromMongoError(err)
	}

	if _, err := s.db.Collection("users").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"name": name}}); err != nil {
		return fromMongoError(err)
	}

	if username != "" && username != user.Username {
		if _, err := s.db.Collection("users").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"username": username}}); err != nil {
			return fromMongoError(err)
		}
	}

	if email != "" && email != user.Email {
		if _, err := s.db.Collection("users").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"email": email}}); err != nil {
			return fromMongoError(err)
		}
	}

	if password != "" {
		if _, err := s.db.Collection("users").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"password": password}}); err != nil {
			return fromMongoError(err)
		}
	}

	return nil
}

func (s *Store) UserCreateToken(ctx context.Context, token *models.UserTokenRecover) error {
	if _, err := primitive.ObjectIDFromHex(token.User); err != nil {
		return err
	}

	if _, err := s.db.Collection("recovery_tokens").InsertOne(ctx, token); err != nil {
		return err
	}

	return nil
}

func (s *Store) UserGetToken(ctx context.Context, id string) (*models.UserTokenRecover, error) {
	token := new(models.UserTokenRecover)
	if err := s.db.Collection("recovery_tokens").FindOne(ctx, bson.M{"user": id}).Decode(&token); err != nil {
		return nil, fromMongoError(err)
	}

	return token, nil
}

func (s *Store) UserDeleteTokens(ctx context.Context, id string) error {
	if _, err := s.db.Collection("recovery_tokens").DeleteMany(ctx, bson.M{"user": id}); err != nil {
		return fromMongoError(err)
	}

	return nil
}

func (s *Store) UserUpdateAccountStatus(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if _, err := s.db.Collection("users").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"confirmed": true}}); err != nil {
		return err
	}

	return nil
}

func (s *Store) UserDelete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fromMongoError(err)
	}

	_, err = s.db.Collection("users").DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fromMongoError(err)
	}

	return nil
}

func (s *Store) UserDetachInfo(ctx context.Context, id string) (map[string][]*models.Namespace, error) {
	findOptions := options.Find()

	cursor, err := s.db.Collection("namespaces").Find(ctx, bson.M{"members": bson.M{"$elemMatch": bson.M{"id": id}}}, findOptions)
	if err != nil {
		return nil, fromMongoError(err)
	}
	defer cursor.Close(ctx)

	namespacesMap := make(map[string][]*models.Namespace, 2)
	ownerNamespaceList := make([]*models.Namespace, 0)
	membersNamespaceList := make([]*models.Namespace, 0)

	for cursor.Next(ctx) {
		namespace := new(models.Namespace)
		if err := cursor.Decode(&namespace); err != nil {
			return nil, fromMongoError(err)
		}

		if namespace.Owner != id {
			membersNamespaceList = append(membersNamespaceList, namespace)
		} else {
			ownerNamespaceList = append(ownerNamespaceList, namespace)
		}
	}

	namespacesMap["member"] = membersNamespaceList
	namespacesMap["owner"] = ownerNamespaceList

	return namespacesMap, nil
}
