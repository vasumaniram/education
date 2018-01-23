//SPDX-License-Identifier: Apache-2.0

var tender = require('./controller.js');

module.exports = function(app){

  app.get('/get_tender/:id', function(req, res){
    tender.get_tender(req, res);
  });
  app.get('/add_tender/:tender', function(req, res){
    tender.add_tender(req, res);
  });
  app.get('/get_all_tenders', function(req, res){
    tender.get_all_tenders(req, res);
  });
  // app.get('/change_holder/:holder', function(req, res){
  //   tuna.change_holder(req, res);
  // });
}
