package main

import (
	"encoding/json"

	"go.etcd.io/bbolt"
)

const profileBacketName string = "Profiles"

type Profile struct {
	Name        string `json:"name"`
	Profile     []byte `json:"profile"`
	AskPassword bool   `json:"ask_password"`
}

func (p *Profile) CreateProfile(db *bbolt.DB) error {

	o := &Profile{
		Name:        p.Name,
		Profile:     p.Profile,
		AskPassword: p.AskPassword,
	}
	data, err := json.Marshal(o)
	if err != nil {
		return err
	}

	db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(profileBacketName))
		if err := b.Put([]byte(o.Name), []byte(data)); err != nil {
			return err
		}
		return nil
	})

	return nil
}

func (p *Profile) ListProfile(db *bbolt.DB) []Profile {

	os := []Profile{}
	o := Profile{}

	db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(profileBacketName))
		b.ForEach(func(k, v []byte) error {
			if err := json.Unmarshal(v, &o); err != nil {
				return err
			}
			os = append(os, o)
			return nil
		})
		return nil
	})
	return os
}

func (p *Profile) DescribeProfile(db *bbolt.DB) *Profile {

	o := Profile{}

	db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(profileBacketName))
		res := b.Get([]byte(p.Name))
		json.Unmarshal(res, &o)
		return nil
	})

	return &o
}

func (p *Profile) DeleteProfile(db *bbolt.DB) error {
	db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(profileBacketName))
		if err := b.Delete([]byte(p.Name)); err != nil {
			return err
		}
		return nil
	})
	return nil
}
