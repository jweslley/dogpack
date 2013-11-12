module.exports = function(grunt){
  grunt.initConfig({
    concat: {
      js: {
        src: [
          'bower_components/jquery/jquery.js',
          'bower_components/handlebars/handlebars.js',
          'bower_components/ember/ember.js',
          'bower_components/moment/moment.js',
          'dogpack.js'
        ],
        dest: 'js/all.js'
      },
      css:{
        src: [
          'bower_components/bootstrap/dist/css/bootstrap.css',
          'dogpack.css'
        ],
        dest: 'css/all.css'
      }
    },
    uglify: {
      dev: {
        options: {
          mangle: false,
          compress: false,
          preserveComments: 'all'
        },
        files: {
          'js/dogpack.js': ['<%= concat.js.dest %>']
        }
      },
      dist: {
        options: {
          mangle: true,
          compress: true
        },
        files: {
          'js/dogpack.js': ['<%= concat.js.dest %>']
        }
      }
    },
    cssmin: {
      css: {
        src: '<%= concat.css.dest %>',
        dest:'css/dogpack.css'
      }
    },
    watch: {
      options: {
        livereload: true,
      },
      js: {
        files: ['dogpack.js'],
        tasks: ['concat:js', 'uglify:dev'],
      },
      css: {
        files: ['dogpack.css'],
        tasks: ['concat:css', 'cssmin'],
      },
      html: {
        files: ['index.html']
      }
    },
  });

  grunt.loadNpmTasks('grunt-contrib-concat');
  grunt.loadNpmTasks('grunt-contrib-uglify');
  grunt.loadNpmTasks('grunt-contrib-cssmin');
  grunt.loadNpmTasks('grunt-contrib-watch');

  grunt.registerTask('default', ['concat', 'uglify:dist', 'cssmin']);
};

