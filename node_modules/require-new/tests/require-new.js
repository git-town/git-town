var expect = require('chai').expect,
    requireNew = require('../src'),
    fixture = __dirname + '/fixtures/rand.js';

describe('require', function () {
    it('second require for the same module will use cached object', function () {
        var a = require(fixture),
            b = require(fixture);

        expect(a).to.equal(b);
    });
});

describe('requireNew', function () {
    beforeEach(function () {
        delete require.cache[require.resolve(fixture)];
    });
    it('requires a new module object', function () {
        var a = requireNew(fixture),
            b = requireNew(fixture);

        expect(a).to.not.equal(b);
    });
    it('use of requireNew before require does not cache the required module', function () {
        var a = requireNew(fixture),
            b = require(fixture);

        expect(a).to.not.equal(b);
    });
    it('use of requireNew after require does not flush the required module cache', function () {
        var a,
            b,
            c;

        a = require(__dirname + '/fixtures/rand.js');
        b = requireNew(__dirname + '/fixtures/rand.js');
        c = require(__dirname + '/fixtures/rand.js');

        expect(a).to.not.equal(b);
        expect(a).to.equal(c);
    });
});

describe('when executed using STDIN', function () {
    it('uses the current working directory to resolve dependency paths', function (cb) {
        var exec = require('child_process').exec,
            command = 'echo "require(\'./src\')(\'./tests/fixtures/hello-world\')" | node';

        exec(command, function (err, stdout) {
            expect(stdout.trim()).to.equal('Hello, World!');

            cb();
        });
    });
});

describe('when executed using eval', function () {
    it('uses the current working directory to resolve dependency paths', function (cb) {
        var exec = require('child_process').exec,
            command = 'node -e "require(\'./src\')(\'./tests/fixtures/hello-world\')"';

        exec(command, function (err, stdout) {
            expect(stdout.trim()).to.equal('Hello, World!');

            cb();
        });
    });
});
