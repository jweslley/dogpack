Dogpack = Ember.Application.create({
  LOG_TRANSITIONS: true,
  ACTION_NAMES: ["", "alert", "restart", "stop", "exec", "unmonitor", "start", "monitor", ""],
  SERVICE_TYPES: ["Filesystem", "Directory", "File", "Process", "Host", "System", "Fifo", "Program"],
  STATUS_NAMES: ["Accessible", "Accessible", "Accessible", "Running", "Online with all services", "Running", "Accessible", "Status ok"],

  getServerStatusCssClass: function(server) {
    var count = 0;
    $.each(server.services, function(i, service) {
      if (service.status == 0) { count++; }
    });

    if (count == server.services.length) {
      return 'status-green';
    } else {
      return 'status-yellow';
    }
  },

  getMonitorStatusDescription: function(status) {
    switch (status) {
      case 0: return "Not monitored";
      case 1: return "Monitored";
      case 2: return "Initializing";
      case 4: return "Waiting";
    }
  },

  getServiceStatus: function(service) {
    if (service.monitor == 0 || service.monitor == 2 || service.monitor == 4) {
      return this.getMonitorStatusDescription(service.monitor);
    } else if (service.status == 0) {
      return this.STATUS_NAMES[service.type];
    } else {
      return service.status_message;
    }
  },

  getServiceStatusCssClass: function(service) {
    if (service.monitor == 0 || service.monitor == 2 || service.monitor == 4) {
      return 'status-gray'
    } else if (service.status == 0) {
      return 'status-green'
    } else {
      return 'status-red';
    }
  }
});

Dogpack.Router.map(function() {
  this.resource('monit', { path: '/servers/:server_name' });
});

Dogpack.IndexRoute = Ember.Route.extend({
  model: function() {
    return $.getJSON('/status').then(function(data) {
      return $.map(data, function(monit) {
        return monit;
      })
    });
  }
});

Dogpack.MonitRoute = Ember.Route.extend({
  setupController: function(controller, monit) {
    controller.set('model', monit);
  },
  model: function(params) {
    return $.getJSON('/status').then(function(data){
      return data[params.server_name];
    });
  }
});

Dogpack.MonitController = Ember.ObjectController.extend({
  actions: {
    unmonitor: function(service){
      this.action(service, 'unmonitor');
    },
    monitor: function(service){
      this.action(service, 'monitor');
    }
  },

  action: function(service, action) {
    var monit = this.get('model');
    $.ajax({
      type: 'POST',
      url: "/action",
      data: {
        server: monit.server.localhostname,
        service: service,
        action: action
      }
    }).done(function(){
      window.location.reload();
    });
  }
});


// helpers

Ember.Handlebars.helper('service-type', function(type) {
  return Dogpack.SERVICE_TYPES[type];
});

Ember.Handlebars.helper('service-status', function(service) {
  return Dogpack.getServiceStatus(service);
});

Ember.Handlebars.helper('server-class', function(server) {
  return Dogpack.getServerStatusCssClass(server);
});

Ember.Handlebars.helper('service-class', function(service) {
  return Dogpack.getServiceStatusCssClass(service);
});

Ember.Handlebars.helper('monitor-status', function(status) {
  return Dogpack.getMonitorStatusDescription(status);
});

Ember.Handlebars.helper('action-name', function(action) {
  return (action == 0) ? '' : 'Pending ' + Dogpack.ACTION_NAMES[action];
});

Ember.Handlebars.helper('format-duration', function(seconds) {
  return moment.duration(seconds, 'seconds').humanize();
});

Ember.Handlebars.helper('format-filesize', function(filesize) {
  filesize *= 1024;
  if (filesize >= 1073741824) {
    filesize = number_format(filesize / 1073741824, 2, '.', '') + ' Gb';
  } else {
    if (filesize >= 1048576) {
      filesize = number_format(filesize / 1048576, 2, '.', '') + ' Mb';
    } else {
      if (filesize >= 1024) {
        filesize = number_format(filesize / 1024, 0) + ' Kb';
      } else {
        filesize = number_format(filesize, 0) + ' bytes';
      };
    };
  };
  return filesize;
});

function number_format(number, decimals, dec_point, thousands_sep) {
  // http://kevin.vanzonneveld.net
  // +   original by: Jonas Raoni Soares Silva (http://www.jsfromhell.com)
  // +   improved by: Kevin van Zonneveld (http://kevin.vanzonneveld.net)
  // +         bugfix by: Michael White (http://crestidg.com)
  // +         bugfix by: Benjamin Lupton
  // +         bugfix by: Allan Jensen (http://www.winternet.no)
  // +        revised by: Jonas Raoni Soares Silva (http://www.jsfromhell.com)
  // *         example 1: number_format(1234.5678, 2, '.', '');
  // *         returns 1: 1234.57

  var n = number, c = isNaN(decimals = Math.abs(decimals)) ? 2 : decimals;
  var d = dec_point == undefined ? "," : dec_point;
  var t = thousands_sep == undefined ? "." : thousands_sep, s = n < 0 ? "-" : "";
  var i = parseInt(n = Math.abs(+n || 0).toFixed(c)) + "", j = (j = i.length) > 3 ? j % 3 : 0;

  return s + (j ? i.substr(0, j) + t : "") + i.substr(j).replace(/(\d{3})(?=\d)/g, "$1" + t) + (c ? d + Math.abs(n - i).toFixed(c).slice(2) : "");
}
