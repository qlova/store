package db

import (
	"github.com/google/uuid"
	"qlova.org/should"
)

var _ Viewable = UUID{}
var _ Variable = &UUID{}
var _ value = &UUID{}

type uid = uuid.UUID

//Handle mocks
var shouldNotError = should.NotError
var shouldError = should.Error
var shouldBe = should.Be
