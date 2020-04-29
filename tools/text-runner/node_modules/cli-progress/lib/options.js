// global options storage
const _options = {};

// utility to merge defaults
function mergeOption(v, defaultValue){
    if (typeof v === 'undefined' || v === null){
        return defaultValue;
    }else{
        return v;
    }
}

module.exports = {
    // set global options
    parse: function parse(rawOptions, preset){

        // merge preset
        const opt = Object.assign({}, preset, rawOptions);

        // the max update rate in fps (redraw will only triggered on value change)
        _options.throttleTime = 1000 / (mergeOption(opt.fps, 10));

        // the output stream to write on
        _options.stream = mergeOption(opt.stream, process.stderr);

        // external terminal provided ?
        _options.terminal = mergeOption(opt.terminal, null);

        // clear on finish ?
        _options.clearOnComplete = mergeOption(opt.clearOnComplete, false);

        // stop on finish ?
        _options.stopOnComplete = mergeOption(opt.stopOnComplete, false);

        // size of the progressbar in chars
        _options.barsize = mergeOption(opt.barsize, 40);

        // position of the progress bar - 'left' (default), 'right' or 'center'
        _options.align = mergeOption(opt.align, 'left');

        // hide the cursor ?
        _options.hideCursor = mergeOption(opt.hideCursor, false);

        // disable linewrapping ?
        _options.linewrap = mergeOption(opt.linewrap, false);

        // pre-render bar strings (performance)
        _options.barCompleteString = (new Array(_options.barsize + 1 ).join(opt.barCompleteChar || '='));
        _options.barIncompleteString = (new Array(_options.barsize + 1 ).join(opt.barIncompleteChar || '-'));

        // glue sequence (control chars) between bar elements ?
        _options.barGlue = mergeOption(opt.barGlue, '');

        // the bar format
        _options.format = mergeOption(opt.format, 'progress [{bar}] {percentage}% | ETA: {eta}s | {value}/{total}');

        // external time-format provided ?
        _options.formatTime = mergeOption(opt.formatTime, null);

        // external value-format provided ?
        _options.formatValue = mergeOption(opt.formatValue, null);

        // external bar-format provided ?
        _options.formatBar = mergeOption(opt.formatBar, null);

        // the number of results to average ETA over
        _options.etaBufferLength = mergeOption(opt.etaBuffer, 10);

        // allow synchronous updates ?
        _options.synchronousUpdate = mergeOption(opt.synchronousUpdate, true);

        // notty mode
        _options.noTTYOutput = mergeOption(opt.noTTYOutput, false);

        // schedule - 2s
        _options.notTTYSchedule = mergeOption(opt.notTTYSchedule, 2000);
        
        // emptyOnZero - false
        _options.emptyOnZero = mergeOption(opt.emptyOnZero, false);

        // force bar redraw even if progress did not change
        _options.forceRedraw = mergeOption(opt.forceRedraw, false);

        // automated padding to fixed width ?
        _options.autopadding = mergeOption(opt.autopadding, false);

        // autopadding character - empty in case autopadding is disabled
        _options.autopaddingChar = _options.autopadding ? mergeOption(opt.autopaddingChar, '   ') : '';

        return _options;
    },

    // fetch all options
    getOptions: function getOptions(){
        return _options;
    }
};