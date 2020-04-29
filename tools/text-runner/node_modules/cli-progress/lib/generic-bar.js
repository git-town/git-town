const _ETA = require('./eta');
const _Terminal = require('./terminal');
const _formatter = require('./formatter');
const _EventEmitter = require('events');

// Progress-Bar constructor
module.exports = class GenericBar extends _EventEmitter{

    constructor(options){
        super();

        // store options
        this.options = options;

        // store terminal instance
        this.terminal = (this.options.terminal) ? this.options.terminal : new _Terminal(this.options.stream);

        // the current bar value
        this.value = 0;

        // the end value of the bar
        this.total = 100;

        // last drawn string - only render on change!
        this.lastDrawnString = null;

        // start time (used for eta calculation)
        this.startTime = null;

        // last update time
        this.lastRedraw = Date.now();

        // default eta calulator (will be re-create on start)
        this.eta = new _ETA(this.options.etaBufferLength, 0, 0);

        // payload data
        this.payload = {};

        // progress bar active ?
        this.isActive = false;

        // use default formatter or custom one ?
        this.formatter = (typeof this.options.format === 'function') ? this.options.format : _formatter;
    }

    // internal render function
    render(){
        // calculate the normalized current progress
        let progress = (this.value/this.total);

        // handle NaN Errors caused by total=0. Set to complete in this case
        if (isNaN(progress)){
            progress = (this.options && this.options.emptyOnZero) ? 0.0 : 1.0;
        }

        // limiter
        progress = Math.min(Math.max(progress, 0.0), 1.0);

        // formatter params
        const params = {
            progress: progress,
            eta: this.eta.getTime(),
            startTime: this.startTime,
            total: this.total,
            value: this.value,
            maxWidth: this.terminal.getWidth()
        };

        // format string
        const s = this.formatter(this.options, params, this.payload);

        const forceRedraw = this.options.forceRedraw
            // force redraw in notty-mode!
            || (this.options.noTTYOutput && !this.terminal.isTTY());

        // string changed ? only trigger redraw on change!
        if (forceRedraw || this.lastDrawnString != s){
            // trigger event
            this.emit('redraw-pre');

            // set cursor to start of line
            this.terminal.cursorTo(0, null);

            // write output
            this.terminal.write(s);

            // clear to the right from cursor
            this.terminal.clearRight();

            // store string
            this.lastDrawnString = s;

            // set last redraw time
            this.lastRedraw = Date.now();

            // trigger event
            this.emit('redraw-post');
        }
    }

    // start the progress bar
    start(total, startValue, payload){
        // set initial values
        this.value = startValue || 0;
        this.total = (typeof total !== 'undefined' && total >= 0) ? total : 100;
        this.payload = payload || {};

        this.startTime = Date.now();
        this.lastDrawnString = '';

        // initialize eta buffer
        this.eta = new _ETA(this.options.etaBufferLength, this.startTime, this.value);

        // set flag
        this.isActive = true;

        // start event
        this.emit('start', total, startValue);
    }

    // stop the bar
    stop(){
        // set flag
        this.isActive = false;

        // stop event
        this.emit('stop', this.total, this.value);
    }

    // update the bar value
    update(current, payload){
        if (typeof current !== 'undefined' && current !== null) {
            // update value
            this.value = current;

            // add new value; recalculate eta
            this.eta.update(Date.now(), current, this.total);
        }

        // update event (before stop() is called)
        this.emit('update', this.total, this.value);

        // merge payload
        const payloadData = payload || {};
        for (const key in payloadData){
            this.payload[key] = payloadData[key];
        }

        // limit reached ? autostop set ?
        if (this.value >= this.getTotal() && this.options.stopOnComplete) {
            this.stop();
        }
    }

    // update the bar value
    increment(step, payload){
        step = step || 1;
        this.update(this.value + step, payload);
    }

    // get the total (limit) value
    getTotal(){
        return this.total;
    }

    // set the total (limit) value
    setTotal(total){
        if (typeof total !== 'undefined' && total >= 0){
            this.total = total;
        }
    }
}
