// helper functions
var SERVICE_TYPES = ["Filesystem", "Directory", "File", "Process", "Host", "System", "Fifo", "Program"];
var STATUS_NAMES = ["Accessible", "Accessible", "Accessible", "Running", "Online with all services", "Running", "Accessible", "Status ok"];
var MONITOR_STATUS = {0: "Not monitored", 1: "Monitored", 2: "Initializing", 4: "Waiting"};

function serviceStatusDesc(service) {
  if (service.monitor == 0 || service.monitor == 2 || service.monitor == 4) {
    return MONITOR_STATUS[service.monitor];

  } else if (service.status == 0) {
    return STATUS_NAMES[service.type];

  } else {
    return service.status_message;
  }
}
function serviceStatusCssClass(service) {
  if (service.monitor == 0 || service.monitor == 2 || service.monitor == 4) {
    return 'gray';

  } else if (service.status == 0) {
    return 'green';

  } else {
    return 'red';

  }
}
function nodeStatusCssClass(monit) {
  var count = 0;
  $.each(monit.services, function(i, service) {
    if (service.status == 0) { count++; }
  });
  if (count == monit.services.length) {
    return 'green';
  } else {
    return 'yellow';
  }
}
function format_filesize(filesize) {
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
}
function number_format(number, decimals, dec_point, thousands_sep) {
  // http://kevin.vanzonneveld.net
  // + original by: Jonas Raoni Soares Silva (http://www.jsfromhell.com)
  // + improved by: Kevin van Zonneveld (http://kevin.vanzonneveld.net)
  // + bugfix by: Michael White (http://crestidg.com)
  // + bugfix by: Benjamin Lupton
  // + bugfix by: Allan Jensen (http://www.winternet.no)
  // + revised by: Jonas Raoni Soares Silva (http://www.jsfromhell.com)
  // * example 1: number_format(1234.5678, 2, '.', '');
  // * returns 1: 1234.57
  var n = number, c = isNaN(decimals = Math.abs(decimals)) ? 2 : decimals;
  var d = dec_point == undefined ? "," : dec_point;
  var t = thousands_sep == undefined ? "." : thousands_sep, s = n < 0 ? "-" : "";
  var i = parseInt(n = Math.abs(+n || 0).toFixed(c)) + "", j = (j = i.length) > 3 ? j % 3 : 0;
  return s + (j ? i.substr(0, j) + t : "") + i.substr(j).replace(/(\d{3})(?=\d)/g, "$1" + t) + (c ? d + Math.abs(n - i).toFixed(c).slice(2) : "");
}


// components

var Dogpack = React.createClass({
  loadStatus: function() {
    $.ajax({
      url: "/status",
      dataType: 'json',
      success: function(data) {
        this.setState({nodes: data});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error("Unable to load status. Error: " + err);
        //this.setState({data: {}});
      }.bind(this)
    });
  },
  handleHashChange: function() {
    window.location.reload();
  },
  getInitialState: function() {
    return { nodes: { } };
  },
  componentDidMount: function() {
    this.loadStatus();
    setInterval(this.loadStatus, this.props.pollInterval);
    window.onhashchange = this.handleHashChange;
  },
  render: function() {
    var nodeId = window.location.hash.substr(1);
    var node = this.state.nodes[nodeId];
    if (node === undefined) {
      return (
        <div>
          <h1>Dogpack</h1>
          <NodeList nodes={this.state.nodes} />
        </div>
      );
    }
    return (
      <div>
        <h1><a href='/'>Dogpack</a> / {nodeId}</h1>
        <ul className="list">
          <Node ip={nodeId} monit={node} />
        </ul>
        <ServiceList type={5} services={node.services} />
        <ServiceList type={3} services={node.services} />
        <ServiceList type={0} services={node.services} />
        <ServiceList type={1} services={node.services} />
        <ServiceList type={2} services={node.services} />
        <ServiceList type={4} services={node.services} />
        <ServiceList type={6} services={node.services} />
        <ServiceList type={7} services={node.services} />
      </div>
    );
  }
});

var NodeList = React.createClass({
  render: function() {
    var nodes = $.map(this.props.nodes, function(instance, ip) {
      return (<Node ip={ip} monit={instance}/>);
    });
    return (
      <ul className="list">
        {nodes}
      </ul>
    );
  }
});

var Node = React.createClass({
  render: function() {
    return (
      <li className={nodeStatusCssClass(this.props.monit)}>
        <a href={'#'+this.props.ip}>
          <ul className="info">
            <li className="title">{this.props.monit.server.localhostname}</li>
            <li>CPU: {this.props.monit.platform.cpu} - Memory: {format_filesize(this.props.monit.platform.memory)} - IP: {this.props.ip}</li>
            <li>{this.props.monit.platform.name} {this.props.monit.platform.machine} - Kernel {this.props.monit.platform.release}</li>
          </ul>
        </a>
      </li>
    );
  }
});

var ServiceList = React.createClass({
  render: function() {
    var count = 0;
    var serviceType = this.props.type;
    var services = this.props.services.map(function(service) {
      if (serviceType != service.type) {
        return '';
      }
      count++;
      return (<Service data={service} />);
    });

    if (count == 0) return (<div className={SERVICE_TYPES[serviceType]}/>);
    return (
      <div className={SERVICE_TYPES[serviceType]}>
        <h3>{SERVICE_TYPES[serviceType]}</h3>
        <ul className="list">
          {services}
        </ul>
      </div>
    );
  }
});

var Service = React.createClass({
  render: function() {
    return (
      <li className={serviceStatusCssClass(this.props.data)}>
        <ul className="info">
          <li className="title">{this.props.data.name}</li>
          <li>{serviceStatusDesc(this.props.data)}</li>
        </ul>
      </li>
    );
  }
});

React.render(
  <Dogpack pollInterval={10000} />,
  document.getElementById('container')
);
